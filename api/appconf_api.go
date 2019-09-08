package api

import (
	"fmt"
	"github.com/zouxinjiang/le/core"
	"github.com/zouxinjiang/le/models"
	"github.com/zouxinjiang/le/pkgs/cerror"
	"strings"
)

type InnerAppConf string

const (
	TwoFactorState     InnerAppConf = "auth.twofactor.state"
	TwoFactorEmail     InnerAppConf = "auth.twofactor.email"
	TwoFactorImageCode InnerAppConf = "auth.twofactor.imagecode"
)

var allAppConf = []InnerAppConf{}

type AppConfApi struct {
	core.Service
}

func (self AppConfApi) Set(key InnerAppConf, value string) error {
	db := self.DbEng()
	if db == nil {
		return cerror.NewJsonError(core.ErrCode_DbConnectFailed)
	}
	if key == "" {
		return nil
	}
	sqlStr := `INSERT INTO appconf(name,value,state) VALUES(?,?,1) ON CONFLICT (name) UPDATE SET value=excluded.value,state=excluded.state`
	db = db.Exec(sqlStr, key, value)
	return db.Error
}

func (self AppConfApi) SetBatch(params map[InnerAppConf]string) error {
	db := self.DbEng()
	if db == nil {
		return cerror.NewJsonError(core.ErrCode_DbConnectFailed)
	}
	if len(params) == 0 {
		return nil
	}
	sqlStr := `INSERT INTO appconf(name,value,state) VALUES ${value} ON CONFLICT (name) UPDATE SET value=excluded.value,state=excluded.state`
	tmp := []string{}
	vals := []interface{}{}
	for k, v := range params {
		vals = append(vals, k, v)
		tmp = append(tmp, fmt.Sprintf("(?,?,1)"))
	}
	sqlStr = strings.ReplaceAll(sqlStr, "${value}", strings.Join(tmp, ","))
	db = db.Exec(sqlStr, vals...)
	return db.Error
}

func (self AppConfApi) Get(key InnerAppConf) (models.AppConfMdl, error) {
	var res = models.AppConfMdl{}
	db := self.DbEng()
	if db == nil {
		return res, cerror.NewJsonError(core.ErrCode_DbConnectFailed)
	}
	sqlStr := `SELECT * FROM appconf WHERE state=1 AND name=?`
	db = db.Raw(sqlStr, key).First(&res)
	return res, db.Error
}

func (self AppConfApi) GetBatch(key ...InnerAppConf) ([]models.AppConfMdl, error) {
	var res = []models.AppConfMdl{}
	db := self.DbEng()
	if db == nil {
		return res, cerror.NewJsonError(core.ErrCode_DbConnectFailed)
	}
	if len(key) == 0 {
		return res, nil
	}
	sqlStr := `SELECT * FROM appconf WHERE state=1 AND name IN (?)`
	db = db.Raw(sqlStr, key).Find(&res)
	return res, db.Error
}

func (self AppConfApi) GetAll() ([]models.AppConfMdl, int, error) {
	var res = []models.AppConfMdl{}
	var cnt = 0
	db := self.DbEng()
	if db == nil {
		return res, cnt, cerror.NewJsonError(core.ErrCode_DbConnectFailed)
	}
	sqlStr := `SELECT * FROM appconf WHERE state=1`
	db = db.Raw(sqlStr).Find(&res)
	return res, len(res), db.Error
}

func (self AppConfApi) FilterInnerAppConf(data []string) (outer []string, inner []InnerAppConf) {
	inner = []InnerAppConf{}
	outer = []string{}
	for _, v := range data {
		var find = false
		for _, in := range allAppConf {
			if in == InnerAppConf(v) {
				inner = append(inner, in)
				find = true
			}
		}
		if !find {
			outer = append(outer, v)
		}
	}
	return
}

func (self AppConfApi) IsInnerAppConf(data string) bool {
	for _, v := range allAppConf {
		if v == InnerAppConf(data) {
			return true
		}
	}
	return false
}
