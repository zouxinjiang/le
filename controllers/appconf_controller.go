package controllers

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/zouxinjiang/le/api"
	"github.com/zouxinjiang/le/core"
	"github.com/zouxinjiang/le/models"
	"github.com/zouxinjiang/le/pkgs/cerror"
	"github.com/zouxinjiang/le/pkgs/clog"
	"github.com/zouxinjiang/le/pkgs/constraint"
)

type AppConfController struct {
	core.Controller
}

func (self AppConfController) GetAppConfig(c echo.Context) error {
	var params = struct {
		IsAll int    `json:"IsAll" form:"IsAll" query:"IsAll" constraint:"Enum:[0,1]"`
		Names string `json:"Names" form:"Names" query:"Names"`
	}{}
	if err := c.Bind(&params); err != nil {
		return cerror.NewJsonError(core.ErrCode_InvalidParams)
	}
	if err := constraint.Valid(&params); err != nil {
		return cerror.NewJsonErrorWithParams(core.ErrCode_InvalidParam, map[string]interface{}{
			"field":  err.Error(),
			"reason": "value must one of [0,1]",
		})
	}
	appconfApi := api.AppConfApi{}
	var err error
	var fdata = map[string]string{}
	var res = []models.AppConfMdl{}
	if params.IsAll == 1 {
		res, _, err = appconfApi.GetAll()
		if err != nil {
			return err
		}
	} else {
		var names = []api.InnerAppConf{}
		if err := json.Unmarshal([]byte(params.Names), &names); err != nil {
			return cerror.NewJsonErrorWithParams(core.ErrCode_InvalidParam, map[string]interface{}{
				"field":  "Names",
				"reason": "value must json array [string]",
			})
		}
		res, err = appconfApi.GetBatch(names...)
	}
	for _, v := range res {
		fdata[v.Name] = v.Value
	}
	return self.RespJson(c, fdata)
}

func (self AppConfController) SetAppConfig(c echo.Context) error {
	var params = struct {
		Name       string `json:"Name" form:"Name" query:"Name"`
		Value      string `json:"Value" form:"Value" query:"Value"`
		NameValues string `json:"NameValues" form:"NameValues" query:"NameValues"`
	}{}
	if err := c.Bind(&params); err != nil {
		clog.Println(clog.Lvl_Info, err)
		return cerror.NewJsonError(core.ErrCode_InvalidParams)
	}
	var tmp = []struct {
		Name  string `json:"Name"`
		Value string `json:"Value"`
	}{}
	appconfApi := api.AppConfApi{}
	var inner = []api.InnerAppConf{}
	var vals = []string{}
	if params.NameValues != "" {
		if err := json.Unmarshal([]byte(params.NameValues), &tmp); err != nil {
			return cerror.NewJsonErrorWithParams(core.ErrCode_InvalidParam, map[string]interface{}{
				"field":  "NameValues",
				"reason": `value must json object string {"Name":"","Value":""}`,
			})
		}
		var names = []string{}
		for _, v := range tmp {
			names = append(names, v.Name)
			if appconfApi.IsInnerAppConf(v.Value) {
				vals = append(vals, v.Value)
			}
		}
		_, inner = appconfApi.FilterInnerAppConf(names)
	} else {
		if appconfApi.IsInnerAppConf(params.Name) {
			inner = append(inner, api.InnerAppConf(params.Name))
			vals = append(vals, params.Value)
		}
	}
	var kv = map[api.InnerAppConf]string{}
	for i, item := range inner {
		kv[item] = vals[i]
	}
	err := appconfApi.SetBatch(kv)
	if err != nil {
		return err
	}
	return self.RespJson(c, nil)
}
