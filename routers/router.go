package routers

import (
	"github.com/labstack/echo"
	"github.com/zouxinjiang/le/controllers"
	"github.com/zouxinjiang/le/modules"
)

func Init(e *echo.Echo) {
	commonCtl := &controllers.CommonController{}
	appConfCtl := &controllers.AppConfController{}

	registerModule(commonCtl, appConfCtl)

	e.GET("/Common/WeChat/ApiCheck", commonCtl.WeiXinApiCheck)
	e.GET("/Config/AppConfig", appConfCtl.GetAppConfig)

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
