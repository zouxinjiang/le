package api

import (
	"fmt"
	"github.com/zouxinjiang/le/config"
	"github.com/zouxinjiang/le/core"
	"github.com/zouxinjiang/le/models"
	"github.com/zouxinjiang/le/pkgs/cache"
	"github.com/zouxinjiang/le/pkgs/cerror"
	"github.com/zouxinjiang/le/pkgs/clog"
	"github.com/zouxinjiang/le/pkgs/lib"
	"github.com/zouxinjiang/le/pkgs/twofactor"
	"time"
)

var (
	ErrCode_AuthFailed = core.CustomErrorCode{
		Code:    "AuthFailed",
		Message: "authentication failed as ${reason}",
		Params:  nil,
	}
	ErrCode_InvalidPassword = core.CustomErrorCode{
		Code:    "InvalidPassword",
		Message: "invalid password",
		Params:  nil,
	}
	ErrCode_InvalidUserName = core.CustomErrorCode{
		Code:    "InvalidUserName",
		Message: "invalid username",
		Params:  nil,
	}
	ErrCode_TwoFactorError = core.CustomErrorCode{
		Code:    "TwoFactorError",
		Message: "two factor do something wrong",
		Params:  nil,
	}
)

type twoFactorItem struct {
	tfType   twofactor.TwoFactorType
	val      string
	username string
}

var (
	tfcache = cache.NewCache("memory")
)

type UserApi struct {
	core.Service
}

func (self UserApi) GetUserByUserName(username string) (models.UserMdl, error) {
	var res = models.UserMdl{}
	db := self.DbEng()
	if db == nil {
		return res, cerror.NewJsonError(core.ErrCode_DbConnectFailed)
	}
	sqlStr := `SELECT * FROM "user" WHERE username=?`
	db = db.Raw(sqlStr, username).First(&res)
	return res, db.Error
}

func (self UserApi) GetUserById(id int64) (models.UserMdl, error) {
	var res = models.UserMdl{}
	db := self.DbEng()
	if db == nil {
		return res, cerror.NewJsonError(core.ErrCode_DbConnectFailed)
	}
	sqlStr := `SELECT * FROM "user" WHERE id=?`
	db = db.Raw(sqlStr, id).First(&res)
	return res, db.Error
}

func (self UserApi) AuthenticationUserPassword(userName, password string) (map[twofactor.TwoFactorType]string, bool, error) {
	var tf = map[twofactor.TwoFactorType]string{}
	var need = false
	user, err := self.GetUserByUserName(userName)
	if err != nil {
		return tf, false, err
	}
	if user.Id == 0 {
		return tf, false, cerror.NewJsonErrorWithParams(core.ErrCode_RecordNotExist, map[string]interface{}{
			"record": fmt.Sprintf("username=%s", userName),
		})
	}

	if lib.Hmac256X(password, config.GetConfig("MemoryConfig.EncryptKey")) != string(user.Password) {
		return tf, false, cerror.NewJsonErrorWithParams(ErrCode_AuthFailed, map[string]interface{}{
			"reason": "invalid password",
		})
	}
	tfparams := map[string]string{}

	tf = self.GetTwoFactor(userName)
	if len(tf) == 0 {
		return tf, false, nil
	}
	successTf := map[twofactor.TwoFactorType]string{}
	var tfodr = []twofactor.TwoFactorType{"email"}
	for _, v := range tfodr {
		token, ok := tf[v]
		if !ok {
			continue
		}
		res, err := twofactor.New(v).Do(tfparams)
		if err != nil {
			clog.Println(clog.Lvl_Error, string(v), " do some thing error:", err)
			continue
		}
		successTf[v] = token
		need = true
		_ = tfcache.Set(token, twoFactorItem{
			tfType:   v,
			val:      res,
			username: userName,
		}, time.Minute*5)
	}
	if len(successTf) == 0 {
		return successTf, true, cerror.NewJsonError(ErrCode_TwoFactorError)
	}
	return tf, need, err
}

func (self UserApi) GetTwoFactor(userName string) map[twofactor.TwoFactorType]string {
	var fdata = map[twofactor.TwoFactorType]string{}
	appconfApi := AppConfApi{}
	res, err := appconfApi.GetBatch(TwoFactorState, TwoFactorEmail, TwoFactorImageCode)
	if err != nil {
		return fdata
	}
	for _, v := range res {
		if InnerAppConf(v.Name) == TwoFactorState {
			if v.Value == "0" {
				fdata = map[twofactor.TwoFactorType]string{}
				break
			}
		}
		if v.Value == "1" {
			fdata[v.Name] = v.Value
		}
	}
	return fdata
}

func (self UserApi) AuthenticationTwoFactor(userName, token string, factor string) bool {
	val := tfcache.Get(token)
	if val == nil {
		return false
	}
	cacheFactor := val.(twoFactorItem)
	if cacheFactor.val == factor {
		return true
	}
	return false
}
