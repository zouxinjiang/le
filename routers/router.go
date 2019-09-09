package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/zouxinjiang/le/controllers"
	"github.com/zouxinjiang/le/middlewares/logincheck"
	"github.com/zouxinjiang/le/modules"
)

func Init(e *echo.Echo) {
	commonCtl := &controllers.CommonController{}
	appConfCtl := &controllers.AppConfController{}
	loginCtl := &controllers.LoginController{}
	userCtl := &controllers.UserController{}

	registerModule(commonCtl, appConfCtl, loginCtl, userCtl)

	e.GET("/Config/AppConfig", appConfCtl.GetAppConfig, logincheck.LoginCheckMiddle)
	e.PATCH("/Config/AppConfig", appConfCtl.SetAppConfig, logincheck.LoginCheckMiddle)

	e.Any("/Common/WeChat/ApiCheck", commonCtl.WeiXinApiCheck)
	e.Any("/User/WeChat/Login", loginCtl.WeChatLogin)
	e.POST("/User/Authentication", loginCtl.AuthenticationUser)
	e.POST("/User/AuthenticationTwoFactor", loginCtl.AuthenticationTwoFactor)
	e.POST("/User/Register", loginCtl.Register)
	e.POST("/User/ForgetPassword", loginCtl.ForgetPassword)
	e.POST("/User/ResetPassword", loginCtl.ResetPassword)

	e.GET("/User/MyInfo", userCtl.GetMyInfo, logincheck.LoginCheckMiddle)
	e.PATCH("/User/MyInfo", userCtl.UpdateUserInfo, logincheck.LoginCheckMiddle)
	e.PATCH("/User/ChangeMyPassword", userCtl.ChangeMyPassword, logincheck.LoginCheckMiddle)

}

var modules_arr = []modules.Module{}

func registerModule(m ...modules.Module) {
	if modules_arr == nil {
		modules_arr = []modules.Module{}
	}
	for _, v := range m {
		_ = v.Init()
		_ = v.Install()
		modules_arr = append(modules_arr, v)
		_ = v.Start(nil)
	}
}

func Free() {
	for _, v := range modules_arr {
		_ = v.Stop(nil)
	}
}
