package gongzhonghao

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	createTime       time.Time
	cacheAccessToken = ApiAccessToken{}
)

type GongZhongHao struct {
	appId   string
	secret  string
	cacheOn bool
	urlPath map[string]string
}

func New(appid, secret string) *GongZhongHao {
	obj := GongZhongHao{
		appId:   appid,
		secret:  secret,
		cacheOn: true,
	}
	obj.initUrl()
	return &obj
}

func (self *GongZhongHao) initUrl() {
	self.urlPath = map[string]string{
		"GetAccessToken":      "https://api.weixin.qq.com/cgi-bin/token",
		"CreateMenu":          "https://api.weixin.qq.com/cgi-bin/menu/create",
		"GetMenu":             "https://api.weixin.qq.com/cgi-bin/menu/get",
		"DeleteMenu":          "https://api.weixin.qq.com/cgi-bin/menu/delete",
		"SetIndustry":         "https://api.weixin.qq.com/cgi-bin/template/api_set_industry",
		"GetIndustry":         "https://api.weixin.qq.com/cgi-bin/template/get_industry",
		"GetTemplate":         "https://api.weixin.qq.com/cgi-bin/template/get_all_private_template",
		"DeleteTemplate":      "https://api.weixin.qq.com/cgi-bin/template/del_private_template",
		"SendTemplateMessage": "https://api.weixin.qq.com/cgi-bin/message/template/send",
	}
}

func (self *GongZhongHao) SetAppIdSecret(appid, secret string) {
	self.appId = appid
	self.secret = secret
	self.cacheOn = true
	self.initUrl()
}

func (self GongZhongHao) GetAccessToken() (ApiAccessToken, error) {
	if self.cacheOn {
		if cacheAccessToken.AccessToken != "" &&
			createTime.Add(time.Second*time.Duration(cacheAccessToken.ExpiresIn-60)).Unix() < time.Now().Unix() {
			return cacheAccessToken, nil
		}
	}

	var res = ApiAccessToken{}
	var apiErr = ApiError{}
	str := fmt.Sprintf("%s?grant_type=client_credential&appid=%s&secret=%s",
		self.urlPath["GetAccessToken"], self.appId, self.secret)
	resp, err := http.Get(str)
	if err != nil {
		return res, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	err = json.Unmarshal(content, &res)
	_ = json.Unmarshal(content, &apiErr)
	if err == nil && res.AccessToken == "" {
		err = errors.New(fmt.Sprint("code:", apiErr.ErrCode, "message:", apiErr.ErrMsg))
	}
	if err == nil {
		createTime = time.Now()
		cacheAccessToken = res
	}
	return res, err
}
