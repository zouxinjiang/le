package gongzhonghao

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/zouxinjiang/le/pkgs/clog"
)

func (self *GongZhongHao) SetIndustry(primaryCode, SecondaryCode string) error {
	token, err := self.GetAccessToken()
	if err != nil {
		return err
	}
	data := fmt.Sprintf(`{"industry_id1":"%s","industry_id2":"%s"}`, primaryCode, SecondaryCode)

	str := fmt.Sprintf("%s?access_token=%s", self.urlPath["SetIndustry"], token.AccessToken)
	_, err = http.Post(str, "application/json", strings.NewReader(data))
	return err
}

func (self *GongZhongHao) GetIndustry(primaryCode, SecondaryCode string) (ApiIndustry, error) {
	var res = ApiIndustry{}
	token, err := self.GetAccessToken()
	if err != nil {
		return res, err
	}
	str := fmt.Sprintf("%s?access_token=%s", self.urlPath["GetIndustry"], token.AccessToken)

	resp, err := http.Get(str)
	if err != nil {
		return res, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	err = json.Unmarshal(content, &res)
	return res, err
}

func (self *GongZhongHao) GetTemplateList() ([]MessageTemplate, error) {
	var res = []MessageTemplate{}
	token, err := self.GetAccessToken()
	if err != nil {
		return res, err
	}
	str := fmt.Sprintf("%s?access_token=%s", self.urlPath["GetTemplate"], token.AccessToken)
	resp, err := http.Get(str)
	if err != nil {
		return res, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	data := map[string][]MessageTemplate{
		"template_list": res,
	}

	err = json.Unmarshal(content, &data)
	return data["template_list"], err
}

func (self *GongZhongHao) DeleteTemplate(id string) error {
	token, err := self.GetAccessToken()
	if err != nil {
		return err
	}
	str := fmt.Sprintf("%s?access_token=%s", self.urlPath["DeleteTemplate"], token.AccessToken)
	data := fmt.Sprintf(`{"template_id" : "%s"}`, id)
	_, err = http.Post(str, "application/json", strings.NewReader(data))
	return err
}

func (self *GongZhongHao) SendTemplateMessage(tmplId string, toUser string, data []TemplateMessageDataItem, callbackUrl string) (string, error) {
	var msgid = ""
	token, err := self.GetAccessToken()
	if err != nil {
		return msgid, err
	}
	str := fmt.Sprintf("%s?access_token=%s", self.urlPath["SendTemplateMessage"], token.AccessToken)
	var (
		msgData = []string{}
		dataVal = ""
	)
	for _, item := range data {
		color := "#173177"
		if item.Color != "" {
			color = item.Color
		}

		msgData = append(msgData,
			fmt.Sprintf(`"%s":{"value":"%s","color":"%s"}`, item.Keyword, item.Value, color))
	}
	dataVal = fmt.Sprintf(`{
		"touser":"%s",
		"template_id":"%s",
		"url":"%s",
		"data":{%s}
	}`, toUser, tmplId, callbackUrl, strings.Join(msgData, ","))

	clog.Debug("消息内容", dataVal)
	resp, err := http.Post(str, "application/json", strings.NewReader(dataVal))
	if err != nil {
		return msgid, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	respData := map[string]interface{}{}
	err = json.Unmarshal(content, &respData)
	if err == nil && fmt.Sprint(respData["errcode"]) == "0" {
		msgid = fmt.Sprint(respData["msgid"])
	}
	if fmt.Sprint(respData["errcode"]) != "0" {
		err = errors.New(fmt.Sprintf("code:%v message:%v", respData["errcode"], respData["errmsg"]))
	}
	return msgid, err
}
