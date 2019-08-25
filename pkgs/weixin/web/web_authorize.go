/*
 * Copyright (c) 2019.
 */

package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"cooker/pkgs/clog"
	"cooker/pkgs/weixin/lib"
)

type WebAuthorize struct {
	appid                 string
	secret                string
	accessTokenUrl        string
	refreshAccessTokenUrl string
	userInfoUrl           string
}

func NewWebAuthorize(appid, secret string) *WebAuthorize {
	obj := WebAuthorize{
		appid:  appid,
		secret: secret,
	}
	obj.initUrl()
	return &obj
}

func (self *WebAuthorize) initUrl() {
	self.accessTokenUrl = "https://api.weixin.qq.com/sns/oauth2/access_token"
	self.refreshAccessTokenUrl = "https://api.weixin.qq.com/sns/oauth2/refresh_token"
	self.userInfoUrl = "https://api.weixin.qq.com/sns/userinfo"
}

func (self *WebAuthorize) SetAppIdSecret(appid, secret string) {
	self.appid = appid
	self.secret = secret
	self.initUrl()
}

func (self *WebAuthorize) ValidateSignature(factors []string, dstSign string) bool {
	sort.Strings(factors)
	tmpStr := strings.Join(factors, "")
	calcSign := lib.Sha1HexString(tmpStr)
	if calcSign == dstSign {
		return true
	} else {
		return false
	}
}

func (self *WebAuthorize) Code2AccessToken(code string) (ApiAccessToken, error) {
	var accessToken = ApiAccessToken{}
	var apiErr = ApiError{}
	urlStr := fmt.Sprintf("%s?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		self.accessTokenUrl, self.appid, self.secret, code)

	resp, err := http.Get(urlStr)
	if err != nil {
		return accessToken, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	clog.Debug(string(content))

	err = json.Unmarshal(content, &accessToken)
	_ = json.Unmarshal(content, &apiErr)
	if err == nil && accessToken.AccessToken == "" {
		err = errors.New(fmt.Sprint("code:", apiErr.ErrorCode, "msg:", apiErr.ErrMsg))
	}
	return accessToken, err
}

func (self *WebAuthorize) RefreshAccessToken(refreshToken ApiAccessToken) (ApiAccessToken, error) {
	var accessToken = ApiAccessToken{}
	var apiErr = ApiError{}
	urlStr := fmt.Sprintf("%s?appid=%s&grant_type=refresh_token&refresh_token=%s",
		self.refreshAccessTokenUrl, self.appid, refreshToken.RefreshToken)

	resp, err := http.Get(urlStr)
	if err != nil {
		return accessToken, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	clog.Debug(string(content))
	err = json.Unmarshal(content, &accessToken)
	_ = json.Unmarshal(content, &apiErr)
	if err == nil && accessToken.AccessToken == "" {
		err = errors.New(fmt.Sprint("code:", apiErr.ErrorCode, "msg:", apiErr.ErrMsg))
	}
	return accessToken, err
}

func (self *WebAuthorize) GetUserInfo(accessToken ApiAccessToken) (ApiUserInfo, error) {
	var userInfo = ApiUserInfo{}
	var apiErr = ApiError{}
	urlStr := fmt.Sprintf("%s?access_token=%s&openid=%s&lang=zh_CN",
		self.userInfoUrl, accessToken.AccessToken, accessToken.OpenId)

	resp, err := http.Get(urlStr)
	if err != nil {
		return userInfo, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	clog.Debug(string(content))
	err = json.Unmarshal(content, &userInfo)
	_ = json.Unmarshal(content, &apiErr)
	if userInfo.OpenId == "" && err == nil {
		err = errors.New(fmt.Sprint("code:", apiErr.ErrorCode, "msg:", apiErr.ErrMsg))
	}
	return userInfo, err
}
