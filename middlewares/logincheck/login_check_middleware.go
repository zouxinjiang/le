package logincheck

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/zouxinjiang/le/config"
	"github.com/zouxinjiang/le/controllers"
	"github.com/zouxinjiang/le/core"
	"github.com/zouxinjiang/le/pkgs/cerror"
	"github.com/zouxinjiang/le/types"
)

var (
	ErrCode_PermissionDeny = core.CustomErrorCode{
		Code:    "PermissionDeny",
		Message: "user need login",
		Params:  nil,
	}
)

func LoginCheckMiddle(hf echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var loginKeyWord = config.GetConfig("MemoryConfig.LoginKey")
		var loginUser = types.UserSession{}
		var err error
		var errFunc = func() error { return nil }
		var haslogin = false
		// token验证
		token := c.QueryParam(loginKeyWord)
		if token != "" {
			loginUser, err = controllers.Token2UserSession(token)
			if err != nil {
				errFunc = func() error {
					return err
				}
			} else {
				haslogin = true
				controllers.FlushTokenExpireTime(token)
			}
		}
		if !haslogin {
			// session 验证
			sess, err := session.Get(loginKeyWord, c)
			if err != nil {
				errFunc = func() error {
					return cerror.NewJsonError(ErrCode_PermissionDeny)
				}
			}
			var ok bool
			loginUser, ok = sess.Values["LoginUser"].(types.UserSession)
			if !ok {
				errFunc = func() error {
					return cerror.NewJsonError(ErrCode_PermissionDeny)
				}
			} else {
				sess.Options = &sessions.Options{
					Path:   "/",
					MaxAge: 60 * 30,
				}
				_ = sess.Save(c.Request(), c.Response())
				haslogin = true
			}
		}
		if !haslogin {
			return errFunc()
		}
		c.Set("LoginUser", loginUser)
		return hf(c)
	}
}
