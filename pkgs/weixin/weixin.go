package weixin

import (
	"cooker/pkgs/weixin/gongzhonghao"
	"cooker/pkgs/weixin/web"
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
