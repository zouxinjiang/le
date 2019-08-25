package gongzhonghao

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (self *GongZhongHao) CreateCustomMenu(button []CustomMenuButton) error {
	var apiErr = ApiError{}
	token, err := self.GetAccessToken()
	if err != nil {
		return err
	}
	str := fmt.Sprintf("%s?access_token=%s", self.urlPath["CreateMenu"], token.AccessToken)
	data := map[string]interface{}{
		"button": button,
	}
	tmp, _ := json.Marshal(data)
	reader := strings.NewReader(string(tmp))
	resp, err := http.Post(str, "application/json", reader)
	if err != nil {
		return err
	}
	content, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	_ = json.Unmarshal(content, &apiErr)
	return errors.New(fmt.Sprint("code:", apiErr.ErrCode, "message:", apiErr.ErrMsg))
}

func (self *GongZhongHao) GetCustomMenu() ([]CustomMenuButton, error) {
	var res = []CustomMenuButton{}
	token, err := self.GetAccessToken()
	if err != nil {
		return res, err
	}
	str := fmt.Sprintf("%s?access_token=%s", self.urlPath["GetMenu"], token.AccessToken)
	resp, err := http.Get(str)
	content, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return res, err
	}
	data := map[string]map[string][]CustomMenuButton{
		"menu": {
			"button": res,
		},
	}
	err = json.Unmarshal(content, &data)
	return data["menu"]["button"], err
}

func (self *GongZhongHao) DeleteCustomMenu() error {
	var apiErr = ApiError{}
	token, err := self.GetAccessToken()
	if err != nil {
		return err
	}
	str := fmt.Sprintf("%s?access_token=%s", self.urlPath["DeleteMenu"], token.AccessToken)
	resp, err := http.Get(str)
	content, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &apiErr)
	if err != nil {
		return err
	}
	if apiErr.ErrCode == 0 {
		return nil
	}
	return errors.New(fmt.Sprint("code:", apiErr.ErrCode, "message:", apiErr.ErrMsg))
}
