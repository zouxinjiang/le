package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/zouxinjiang/le/config"
	"github.com/zouxinjiang/le/pkgs/clog"
)

var dbInstance *gorm.DB

func databaseString() string {
	ds := ""
	ds = fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		config.GetConfig("FileConfig.DatabaseConfig.Host"),
		config.GetConfig("FileConfig.DatabaseConfig.Port"),
		config.GetConfig("FileConfig.DatabaseConfig.Database"),
		config.GetConfig("FileConfig.DatabaseConfig.UserName"),
		config.GetConfig("FileConfig.DatabaseConfig.Password"),
	)
	clog.Println(clog.Lvl_Debug, ds)
	return ds
}

func ConnectPg(ds ...string) (*gorm.DB, error) {
	if len(ds) > 0 {
		clog.Println(clog.Lvl_Debug, ds)
		return gorm.Open("postgres", ds[0])
	}
	return gorm.Open("postgres", databaseString())
}

func OneInstance() *gorm.DB {
	var idx = 1
RECONNECT:
	var db *gorm.DB
	var tmpdb *gorm.DB
	var err error
	if dbInstance != nil {
		db = dbInstance
	} else {
		tmpdb, err = ConnectPg(databaseString())
		if err != nil {
			dbInstance = nil
			if idx > 3 {
				return nil
			}
			idx++
			goto RECONNECT
		} else {
			dbInstance = tmpdb
		}
		db = tmpdb
	}
	// 检测数据库是否正常
	var IsAliveSql = `SELECT 1+2 AS alive`
	tmpdb = db.Exec(IsAliveSql)
	if tmpdb.Error != nil {
		if idx > 3 {
			return nil
		}
		idx++
		goto RECONNECT
	}
	return db.New()
}
