package main

import (
	"github.com/labstack/echo"
	"github.com/zouxinjiang/le/config"
	"github.com/zouxinjiang/le/pkgs/clog"
	"github.com/zouxinjiang/le/routers"
	"strconv"
)

func init() {
	_ = config.Init()
	show := config.GetConfig("FileConfig.LogConfig.ShowLevel")
	level, _ := strconv.ParseInt(show, 10, 64)
	clog.SetShowLevel(clog.LogLevel(level))
}

func main() {
	e := echo.New()

	routers.Init(e)

	e.HTTPErrorHandler = AppErrorHandleFunc

	addr := config.GetConfig("FileConfig.WebConfig.Address") + ":" + config.GetConfig("FileConfig.WebConfig.Port")
	clog.Info(e.Start(addr))
}
