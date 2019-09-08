package core

import (
	"github.com/jinzhu/gorm"
	"github.com/zouxinjiang/le/config"
	"github.com/zouxinjiang/le/core/db"
)

type Service struct {
}

func (s Service) Init() error {
	_ = s.DbEng()
	return nil
}

func (s Service) Install() error {
	return nil
}

func (s Service) Start(params map[string]string) error {
	return nil
}

func (s Service) Stop(params map[string]string) error {
	return nil
}

//=================================公共部分==========================

func (Service) DbEng() *gorm.DB {
	res := db.OneInstance()
	if res == nil {
		return res
	}
	if config.GetConfig("FileConfig.DatabaseConfig.Debug") == "true" {
		res = res.Debug()
	} else {
		res = res.LogMode(false)
	}
	return res
}

func (Service) NewDbEngInstanceForce() (*gorm.DB, error) {
	res, err := db.ConnectPg()
	if err != nil {
		return nil, err
	}
	if config.GetConfig("FileConfig.DatabaseConfig.Debug") == "true" {
		res = res.Debug()
	}
	return res, err
}

func (self Service) WrapDbErrorCode(err error) CustomErrorCode {
	if IsDbErrorRecordNotFount(err) {
		return ErrCode_RecordNotExist
	} else if IsDbErrorUnique(err) {
		return ErrCode_RecordExisted
	} else if IsDbErrorForeignKey(err) {
		return ErrCode_RecordExisted
	} else {
		return ErrCode_Unknown
	}
}
