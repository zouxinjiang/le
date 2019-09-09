package controllers

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/zouxinjiang/le/api"
	"github.com/zouxinjiang/le/config"
	"github.com/zouxinjiang/le/core"
	"github.com/zouxinjiang/le/pkgs/cache"
	"github.com/zouxinjiang/le/pkgs/cerror"
	"github.com/zouxinjiang/le/pkgs/constraint"
	"github.com/zouxinjiang/le/pkgs/lib"
	"github.com/zouxinjiang/le/pkgs/weixin"
	"github.com/zouxinjiang/le/types"
	"strings"
	"time"
)

type LoginController struct {
	core.Controller
}

var (
	ErrCode_TokenExpired = core.CustomErrorCode{
		Code:    "TokenExpired",
		Message: "token has expired",
		Params:  nil,
	}

	TokenCache          = cache.NewCache("memory")
	userinfoCache       = cache.NewCache("memory")
	tokenExpireInterval = 60 * 30
)

func (self LoginController) SetLoginInfo(c echo.Context, userinfo types.UserSession, timeout int, token string) error {
	var loginKeyWord = config.GetConfig("MemoryConfig.LoginKey")
	//写session,写token
	sess, _ := session.Get(loginKeyWord, c)
	sess.Options = &sessions.Options{
		Path:   "/",
		MaxAge: timeout,
	}
	sess.Values["LoginUser"] = userinfo
	_ = sess.Save(c.Request(), c.Response())
	_ = TokenCache.Set(token, userinfo, time.Duration(timeout)*time.Second)
	return nil
}

func Token2UserSession(token string) (types.UserSession, error) {
	val := TokenCache.Get(token)
	loginUser, ok := val.(types.UserSession)
	if ok {
		return loginUser, nil
	}
	return loginUser, cerror.NewJsonError(ErrCode_TokenExpired)
}

func FlushTokenExpireTime(token string) {
	val := TokenCache.Get(token)
	if _, ok := val.(types.UserSession); ok {
		_ = TokenCache.Set(token, val, time.Duration(tokenExpireInterval)*time.Second)
	}
}

func (self LoginController) AuthenticationUser(c echo.Context) error {
	var params = struct {
		UserName string `json:"UserName" form:"UserName" query:"UserName"`
		Password string `json:"Password" form:"Password" query:"Password"`
	}{}
	if err := c.Bind(&params); err != nil {
		return cerror.NewJsonError(core.ErrCode_InvalidParams)
	}
	if params.UserName == "" {
		return cerror.NewJsonErrorWithParams(core.ErrCode_InvalidParam, map[string]interface{}{
			"field":  "UserName",
			"reason": "value must specific",
		})
	}
	userApi := api.UserApi{}
	userInfo, err := userApi.GetUserByUserName(params.UserName)
	if err != nil {
		return cerror.NewJsonErrorWithParams(self.WrapDbErrorCode(err), map[string]interface{}{
			"record": params.UserName,
		})
	}
	if userInfo.Id == 0 {
		return cerror.NewJsonErrorWithParams(core.ErrCode_RecordNotExist, map[string]interface{}{
			"record": fmt.Sprintf(" user username=%s ", params.UserName),
		})
	}
	tf, _, err := userApi.AuthenticationUserPassword(params.UserName, params.Password)
	if err != nil {
		return err
	}
	loginUser := types.UserSession{
		UserName:  userInfo.Username,
		LoginTime: time.Now(),
	}
	var loginTimeout = tokenExpireInterval
	if len(tf) == 0 {
		var token = lib.RandStr(36)
		_ = self.SetLoginInfo(c, loginUser, loginTimeout, token)
		return self.RespJson(c, map[string]interface{}{
			"Token": token,
		})
	}
	for _, v := range tf {
		_ = userinfoCache.Set(fmt.Sprint(v["Token"]), loginUser, time.Minute*5)
	}
	return self.RespJson(c, map[string]interface{}{
		"TwoFactor": tf,
	})
}

func (self LoginController) AuthenticationTwoFactor(c echo.Context) error {
	var params = struct {
		UserName string `json:"UserName" form:"UserName" query:"UserName"`
		Token    string `json:"Token" form:"Token" query:"Token" constraint:"require"`
		Factor   string `json:"Factor" form:"Factor" query:"Factor" constraint:"require"`
	}{}
	if err := c.Bind(&params); err != nil {
		return cerror.NewJsonError(core.ErrCode_InvalidParams)
	}
	if err := constraint.Valid(params); err != nil {
		return cerror.NewJsonErrorWithParams(core.ErrCode_InvalidParam, map[string]interface{}{
			"field":  err.Error(),
			"reason": "value is required",
		})
	}
	userApi := api.UserApi{}
	if userApi.AuthenticationTwoFactor(params.UserName, params.Token, params.Factor) {
		val := userinfoCache.Get(params.Token)
		loginUser, _ := val.(types.UserSession)
		var loginTimeout = tokenExpireInterval
		var token = lib.RandStr(36)
		_ = self.SetLoginInfo(c, loginUser, loginTimeout, token)
		return self.RespJson(c, map[string]interface{}{
			"Token": token,
		})
	}
	return cerror.NewJsonError(api.ErrCode_TwoFactorWrong)
}

func (self LoginController) Register(c echo.Context) error {
	var params = struct {
		UserName string `json:"UserName" form:"UserName" query:"UserName" constraint:"required"`
		Password string `json:"Password" form:"Password" query:"Password" constraint:"required"`
		Name     string `json:"Name" form:"Name" query:"Name"`
		Email    string `json:"Email" form:"Email" query:"Email" constraint:"type:email"`
		Mobile   string `json:"Mobile" form:"Mobile" query:"Mobile" constraint:"type:mobile"`
	}{}
	if err := c.Bind(&params); err != nil {
		return cerror.NewJsonError(core.ErrCode_InvalidParams)
	}
	if err := constraint.Valid(&params); err != nil {
		return cerror.NewJsonErrorWithParams(core.ErrCode_InvalidParam, map[string]interface{}{
			"field":  err.Error(),
			"reason": "value is not right",
		})
	}
	userApi := api.UserApi{}
	err := userApi.AddUser(params.UserName, params.Name, params.Password, params.Email, params.Mobile)
	if err != nil {
		return err
	}
	return self.RespJson(c, "")
}

func (self LoginController) ForgetPassword(c echo.Context) error {
	var params = struct {
		UserName string `json:"UserName" form:"UserName" query:"UserName"`
		Email    string `json:"Email" form:"Email" query:"Email" constraint:"type:email;required"`
	}{}
	if err := c.Bind(&params); err != nil {
		return cerror.NewJsonError(core.ErrCode_InvalidParams)
	}
	if err := constraint.Valid(&params); err != nil {
		return cerror.NewJsonErrorWithParams(core.ErrCode_InvalidParam, map[string]interface{}{
			"field":  err.Error(),
			"reason": "value is required or invalid",
		})
	}
	userApi := api.UserApi{}
	user, err := userApi.GetUserByUserName(params.UserName)
	if core.IsDbErrorRecordNotFount(err) {
		return cerror.NewJsonErrorWithParams(core.ErrCode_RecordNotExist, map[string]interface{}{
			"record": fmt.Sprintf("user username=%s", params.UserName),
		})
	}
	if user.Email != "" && user.Email == params.Email {
		// 邮箱验证通过，将验证码发到邮箱
		fdata, err := userApi.SendTwoFactor("email", map[string]string{
			"UserName": params.UserName,
			"To":       params.Email,
		})
		if err != nil {
			return cerror.NewJsonError(api.ErrCode_TwoFactorError)
		}
		return self.RespJson(c, fdata)
	}
	return cerror.NewJsonErrorWithParams(core.ErrCode_InvalidParam, map[string]interface{}{
		"field":  "Email",
		"reason": "value is required",
	})
}

func (self LoginController) ResetPassword(c echo.Context) error {
	var params = struct {
		Token       string `json:"Token" form:"Token" query:"Token" constraint:"required"`
		VerifyCode  string `json:"VerifyCode" form:"VerifyCode" query:"VerifyCode" constraint:"required"`
		NewPassword string `json:"NewPassword" form:"NewPassword" query:"NewPassword" constraint:"required"`
	}{}
	if err := c.Bind(&params); err != nil {
		return cerror.NewJsonError(core.ErrCode_InvalidParams)
	}
	if err := constraint.Valid(&params); err != nil {
		return cerror.NewJsonErrorWithParams(core.ErrCode_InvalidParam, map[string]interface{}{
			"field":  err.Error(),
			"reason": "value is required",
		})
	}
	userApi := api.UserApi{}
	info, err := userApi.ValidateTwoFactor(params.Token, params.VerifyCode)
	user, _ := userApi.GetUserByUserName(info.Username)
	if err != nil {
		return err
	}
	err = userApi.ResetUserPassword(user.Id, params.NewPassword)
	if err != nil {
		return err
	}
	return self.RespJson(c, "")
}

func (self LoginController) WeChatLogin(c echo.Context) error {
	var params = struct {
		FailUrl    string `json:"FailUrl" form:"FailUrl" query:"FailUrl"`
		SuccessUrl string `json:"SuccessUrl" form:"SuccessUrl" query:"SuccessUrl"`
		Code       string `json:"Code" form:"Code" query:"Code"`
		State      string `json:"State" form:"State" query:"State"`
	}{}
	if err := c.Bind(&params); err != nil {
		return cerror.NewJsonError(core.ErrCode_InvalidParams)
	}
	wx := weixin.New(config.GetConfig("FileConfig.WeiXinConfig.AppId"), config.GetConfig("FileConfig.WeiXinConfig.Secret"))
	ak, err := wx.WebAuthorize.Code2AccessToken(params.Code)
	if err != nil {
		if strings.Contains(params.FailUrl, "?") {
			params.SuccessUrl += "&ErrorMessage=" + err.Error()
		} else {
			params.FailUrl += "?ErrorMessage=" + err.Error()
		}
		return c.Redirect(302, params.FailUrl)
	}
	userinfo, err := wx.GetUserInfo(ak)
	if err != nil {
		if strings.Contains(params.FailUrl, "?") {
			params.FailUrl += "&ErrorMessage=" + err.Error()
		} else {
			params.FailUrl += "?ErrorMessage=" + err.Error()
		}
		return c.Redirect(302, params.FailUrl)
	}
	userApi := api.UserApi{}
	//查找用户是否存在
	user, _ := userApi.GetUserByUuid(userinfo.OpenId)
	if user.Id == 0 {
		//注册，然后返回登陆信息
		uname := lib.RandStr(10)
		err := userApi.AddWeiXinUser("wx_"+uname, userinfo.OpenId, userinfo.NickName, "123456", "", "")
		if err != nil {
			if strings.Contains(params.FailUrl, "?") {
				params.FailUrl += "&ErrorMessage=" + err.Error()
			} else {
				params.FailUrl += "?ErrorMessage=" + err.Error()
			}
			return c.Redirect(302, params.FailUrl)
		}

		user, _ = userApi.GetUserByUuid(userinfo.OpenId)
	}
	if user.Id == 0 {
		if strings.Contains(params.FailUrl, "?") {
			params.FailUrl += "&ErrorMessage=" + "user login failed"
		} else {
			params.FailUrl += "?ErrorMessage=" + "user login failed"
		}
		return c.Redirect(302, params.FailUrl)
	}
	// 更新头像地址
	_ = userApi.UpdateUserIcon(user.Id, userinfo.HeadImgUrl)
	//产生token，生成session
	var token = lib.RandStr(36)
	loginUser := types.UserSession{
		UserName:  user.Username,
		LoginTime: time.Now(),
	}
	var loginTimeout = tokenExpireInterval
	_ = self.SetLoginInfo(c, loginUser, loginTimeout, token)
	if strings.Contains(params.SuccessUrl, "?") {
		params.SuccessUrl += "&Token=" + token
	} else {
		params.SuccessUrl += "?Token=" + token
	}
	return c.Redirect(302, params.SuccessUrl)
}
