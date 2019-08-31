package api

import (
	"fmt"
	"github.com/zouxinjiang/le/core"
	"github.com/zouxinjiang/le/models"
	"github.com/zouxinjiang/le/pkgs/cerror"
	"strings"
)

type InnerAppConf string

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
	sqlStr := `INSERT INTO appconf(name,value,state) VALUES(?,?,1)`
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
	sqlStr := `INSERT INTO appconf(name,value,state) VALUES`
	tmp := []string{}
	vals := []interface{}{}
	for k, v := range params {
		vals = append(vals, k, v)
		tmp = append(tmp, fmt.Sprintf("(?,?,1)"))
	}
	sqlStr = sqlStr + strings.Join(tmp, ",")
	db = db.Exec(sqlStr, vals...)
	return db.Error
}

func (self AppConfApi) Get(key InnerAppConf) (models.AppConfMdl, error) {
	var res = models.AppConfMdl{}
	db := self.DbEng()
	if db == nil {
		return res, cerror.NewJsonError(core.ErrCode_DbConnectFailed)
	}
	sqlStr := `SELECT * FROM appconf WHERE name=?`
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
	sqlStr := `SELECT * FROM appconf WHERE name IN (?)`
	db = db.Raw(sqlStr, key).First(&res)
	return res, db.Error
}

func (self AppConfApi) GetAll() ([]models.AppConfMdl, int, error) {
	var res = []models.AppConfMdl{}
	var cnt = 0
	db := self.DbEng()
	if db == nil {
		return res, cnt, cerror.NewJsonError(core.ErrCode_DbConnectFailed)
	}
	sqlStr := `SELECT * FROM appconf`
	db = db.Raw(sqlStr).Find(&res)
	return res, len(res), db.Error
}
