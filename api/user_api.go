package api

import (
	"bytes"
	"fmt"
	"github.com/zouxinjiang/le/config"
	"github.com/zouxinjiang/le/core"
	"github.com/zouxinjiang/le/models"
	"github.com/zouxinjiang/le/pkgs/cache"
	"github.com/zouxinjiang/le/pkgs/cerror"
	"github.com/zouxinjiang/le/pkgs/clog"
	"github.com/zouxinjiang/le/pkgs/lib"
	"github.com/zouxinjiang/le/pkgs/twofactor"
	"strings"
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

	ErrCode_InvalidOldPassword = core.CustomErrorCode{
		Code:    "InvalidOldPassword",
		Message: "invalid old password",
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
	ErrCode_TwoFactorWrong = core.CustomErrorCode{
		Code:    "TwoFactorWrong",
		Message: "two factor is not right",
		Params:  nil,
	}
)

type twoFactorItem struct {
	TfType   twofactor.TwoFactorType
	Val      string
	Username string
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

func (self UserApi) GetUserByUuid(uuid string) (models.UserMdl, error) {
	var res = models.UserMdl{}
	db := self.DbEng()
	if db == nil {
		return res, cerror.NewJsonError(core.ErrCode_DbConnectFailed)
	}
	sqlStr := `SELECT * FROM "user" WHERE uuid=?`
	db = db.Raw(sqlStr, uuid).First(&res)
	return res, db.Error
}

func (self UserApi) AuthenticationUserPassword(userName, password string) (map[twofactor.TwoFactorType]map[string]interface{}, bool, error) {
	var fdatatf = map[twofactor.TwoFactorType]map[string]interface{}{}
	var need = false
	user, err := self.GetUserByUserName(userName)
	if err != nil {
		return fdatatf, false, err
	}
	if user.Id == 0 {
		return fdatatf, false, cerror.NewJsonErrorWithParams(core.ErrCode_RecordNotExist, map[string]interface{}{
			"record": fmt.Sprintf("username=%s", userName),
		})
	}

	if lib.Hmac256X(password, config.GetConfig("MemoryConfig.EncryptKey")) != string(user.Password) {
		return fdatatf, false, cerror.NewJsonErrorWithParams(ErrCode_AuthFailed, map[string]interface{}{
			"reason": "invalid password",
		})
	}
	tfparams := map[string]string{
		"UserName": config.GetConfig("FileConfig.EmailConfig.UserName"),
		"Host":     config.GetConfig("FileConfig.EmailConfig.Host"),
		"Port":     config.GetConfig("FileConfig.EmailConfig.Port"),
		"Password": config.GetConfig("FileConfig.EmailConfig.Password"),
		"from":     "[LE]",
		"to":       user.Email,
	}

	tf := self.GetTwoFactor(userName)
	if len(tf) == 0 {
		return fdatatf, false, nil
	}
	var tfodr = []twofactor.TwoFactorType{"email", "imagecode"}
	for _, v := range tfodr {
		_, ok := tf[v]
		if !ok {
			clog.Debug(tf, v)
			continue
		}
		var code = lib.RandNumberStr(6)
		tfparams["code"] = code
		addr, res, err := twofactor.New(v).Do(tfparams)
		if err != nil {
			clog.Error(v, " something went wrong:", err)
			continue
		}
		token := lib.RandStr(32)
		// 写入返回值结果
		fdatatf[v] = map[string]interface{}{
			"Token":   token,
			"Address": addr,
		}
		// 记录code，后续认证使用到
		tf[v] = addr
		need = true
		_ = tfcache.Set(token, twoFactorItem{
			TfType:   v,
			Val:      res,
			Username: userName,
		}, time.Minute*5)
	}
	if len(fdatatf) == 0 {
		return fdatatf, true, cerror.NewJsonError(ErrCode_TwoFactorError)
	}
	return fdatatf, need, err
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
			continue
		}
		if v.Value == "1" {
			if tmp := strings.Split(v.Name, "."); len(tmp) >= 3 {
				fdata[tmp[2]] = v.Value
			}
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
	if cacheFactor.Val == factor {
		return true
	}
	return false
}

func (self UserApi) AddUser(username, name, password, email, mobile string) error {
	db := self.DbEng()
	if db == nil {
		return cerror.NewJsonError(core.ErrCode_DbConnectFailed)
	}
	u, _ := self.GetUserByUserName(username)
	if u.Id != 0 {
		return cerror.NewJsonErrorWithParams(core.ErrCode_RecordExisted, map[string]interface{}{
			"record": fmt.Sprintf(" user username=%s ", username),
		})
	}
	sqlStr := `INSERT INTO "user"(type,username,name,pwd,mobile,email,state,createat,updateat) VALUES (?,?,?,?,?,?,?,?,?)`
	pwd := lib.Hmac256X(password, config.GetConfig("MemoryConfig.EncryptKey"))
	db = db.Exec(sqlStr, "local", username, name, pwd, mobile, email, 1, time.Now(), time.Now())
	return db.Error
}

func (self UserApi) AddWeiXinUser(username, uuid, name, password, email, mobile string) error {
	db := self.DbEng()
	if db == nil {
		return cerror.NewJsonError(core.ErrCode_DbConnectFailed)
	}
	u, _ := self.GetUserByUserName(username)
	if u.Id != 0 {
		return cerror.NewJsonErrorWithParams(core.ErrCode_RecordExisted, map[string]interface{}{
			"record": fmt.Sprintf(" user username=%s ", username),
		})
	}
	sqlStr := `INSERT INTO "user"(type,uuid,username,name,pwd,mobile,email,state,createat,updateat) VALUES (?,?,?,?,?,?,?,?,?,?)`
	pwd := lib.Hmac256X(password, config.GetConfig("MemoryConfig.EncryptKey"))
	db = db.Exec(sqlStr, "wechat", uuid, username, name, pwd, mobile, email, 1, time.Now(), time.Now())
	return db.Error
}

func (self UserApi) UpdateUser(uid int64, info map[string]interface{}) error {
	var keyMap = map[string]string{
		"Name":     "name",
		"Mobile":   "mobile",
		"Email":    "email",
		"UpdateAt": "updateat",
	}

	db := self.DbEng()
	if db == nil {
		return cerror.NewJsonError(core.ErrCode_DbConnectFailed)
	}
	sqlStr := `UPDATE "user" SET ${field} WHERE id=?`
	var fields = []string{}
	var vals = []interface{}{}
	if len(info) > 0 {
		info["UpdateAt"] = time.Now()
	}

	for k, v := range info {
		f1, ok := keyMap[k]
		if !ok {
			continue
		}
		fields = append(fields, fmt.Sprintf("%s=?", f1))
		vals = append(vals, v)
	}
	if len(fields) == 0 {
		return nil
	}
	sqlStr = strings.ReplaceAll(sqlStr, "${field}", strings.Join(fields, ","))
	vals = append(vals, uid)
	db = db.Exec(sqlStr, vals...)
	return db.Error
}

func (self UserApi) UpdateUserIcon(uid int64, url string) error {
	db := self.DbEng()
	if db == nil {
		return cerror.NewJsonError(core.ErrCode_DbConnectFailed)
	}
	sqlStr := `UPDATE "user" SET icon=? WHERE id=?`
	db = db.Exec(sqlStr, url, uid)
	return db.Error
}

func (self UserApi) ChangeUserPassword(uid int64, oldPwd, newPwd string) error {
	db := self.DbEng()
	if db == nil {
		return cerror.NewJsonError(core.ErrCode_DbConnectFailed)
	}
	u, err := self.GetUserById(uid)
	if core.IsDbErrorRecordNotFount(err) || u.Id == 0 {
		return cerror.NewJsonErrorWithParams(core.ErrCode_RecordNotExist, map[string]interface{}{
			"record": fmt.Sprintf(" user id=%d", uid),
		})
	}
	old := lib.Hmac256X(oldPwd, config.GetConfig("MemoryConfig.EncryptKey"))
	if !bytes.Equal(u.Password, []byte(old)) {
		return cerror.NewJsonError(ErrCode_InvalidOldPassword)
	}

	newPassword := lib.Hmac256X(newPwd, config.GetConfig("MemoryConfig.EncryptKey"))
	sqlStr := `UPDATE "user" SET pwd=?,updateat=? WHERE id=?`
	db = db.Exec(sqlStr, newPassword, time.Now(), uid)
	return db.Error
}

func (self UserApi) SendTwoFactor(tpy twofactor.TwoFactorType, params map[string]string) (fdata map[twofactor.TwoFactorType]map[string]string, err error) {
	username := params["UserName"]
	to := params["To"]

	tfparams := map[string]string{
		"UserName": config.GetConfig("FileConfig.EmailConfig.UserName"),
		"Host":     config.GetConfig("FileConfig.EmailConfig.Host"),
		"Port":     config.GetConfig("FileConfig.EmailConfig.Port"),
		"Password": config.GetConfig("FileConfig.EmailConfig.Password"),
		"from":     "[LE]",
		"to":       to,
	}
	addr, res, err := twofactor.New(tpy).Do(tfparams)
	if err != nil {
		clog.Error(" something went wrong:", err)
		return nil, err
	}
	token := lib.RandStr(32)
	// 写入返回值结果
	fdata[tpy] = map[string]string{
		"Token":   token,
		"Address": addr,
	}
	// 记录code，后续认证使用到
	_ = tfcache.Set(token, twoFactorItem{
		TfType:   tpy,
		Val:      res,
		Username: username,
	}, time.Minute*5)
	return fdata, nil
}

func (self UserApi) ValidateTwoFactor(token, factor string) (twoFactorItem, error) {
	res := twoFactorItem{}
	val := tfcache.Get(token)
	if val == nil {
		return res, cerror.NewJsonError(ErrCode_TwoFactorWrong)
	}
	cacheFactor := val.(twoFactorItem)
	if cacheFactor.Val == factor {
		return res, nil
	}
	res = cacheFactor
	return res, cerror.NewJsonError(ErrCode_TwoFactorWrong)
}

func (self UserApi) ResetUserPassword(uid int64, newPassword string) error {
	db := self.DbEng()
	if db == nil {
		return cerror.NewJsonError(core.ErrCode_DbConnectFailed)
	}
	newpwd := lib.Hmac256X(newPassword, config.GetConfig("MemoryConfig.EncryptKey"))
	sqlStr := `UPDATE "user" SET pwd=?,updateat=? WHERE id=?`
	db = db.Exec(sqlStr, newpwd, time.Now(), uid)
	return db.Error
}
