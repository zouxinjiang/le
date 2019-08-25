package main

import (
	"github.com/labstack/echo"
	"github.com/zouxinjiang/le/config"
	"github.com/zouxinjiang/le/pkgs/clog"
	"github.com/zouxinjiang/le/routers"
)

func init() {
	_ = config.Init()
}

func main() {
	e := echo.New()

	routers.Init(e)

	e.HTTPErrorHandler = AppErrorHandleFunc

	addr := config.GetConfig("FileConfig.WebConfig.Address") + ":" + config.GetConfig("FileConfig.WebConfig.Port")
	clog.Info(e.Start(addr))
}
