package controllers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/zouxinjiang/le/config"
	"github.com/zouxinjiang/le/core"
	"github.com/zouxinjiang/le/pkgs/cerror"
	"github.com/zouxinjiang/le/pkgs/weixin/web"
)

type CommonController struct {
	core.Controller
}

func (self CommonController) WeiXinApiCheck(c echo.Context) error {
	webAuthApi := web.NewWebAuthorize(config.GetConfig("FileConfig.WeiXinConfig.AppId"), config.GetConfig("FileConfig.WeiXinConfig.Secret"))

	var params = struct {
		Signature string `json:"signature" form:"signature" query:"signature"`
		Timestamp int64  `json:"timestamp" form:"timestamp" query:"timestamp"`
		Nonce     int64  `json:"nonce" form:"nonce" query:"nonce"`
		EchoStr   string `json:"echostr" form:"echostr" query:"echostr"`
	}{}
	if err := c.Bind(&params); err != nil {
		return cerror.NewJsonError(core.ErrCode_InvalidParams)
	}

	store := []string{
		fmt.Sprintf("%d", params.Timestamp),
		fmt.Sprintf("%d", params.Nonce),
		config.GetConfig("FileConfig.WeiXinConfig.Token"),
	}

	if webAuthApi.ValidateSignature(store, params.Signature) {
		return c.String(200, params.EchoStr)
	} else {
		return c.String(200, "auth failed")
	}
}
