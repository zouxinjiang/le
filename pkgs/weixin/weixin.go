package weixin

import (
	"github.com/zouxinjiang/le/pkgs/weixin/gongzhonghao"
	"github.com/zouxinjiang/le/pkgs/weixin/web"
)

type WeiXin struct {
	*web.WebAuthorize
	*gongzhonghao.GongZhongHao
}

func New(appid, secret string) *WeiXin {
	obj := WeiXin{
		WebAuthorize: web.NewWebAuthorize(appid, secret),
		GongZhongHao: gongzhonghao.New(appid, secret),
	}
	return &obj
}
