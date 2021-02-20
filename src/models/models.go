package models

import (
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"myProject/common/response"
	"myProject/conf"
	"myProject/lib/log"
)

var once sync.Once

var dbInstance *gorm.DB

func Open(config *conf.DBConfig) error {
	var err error
	once.Do(func() {
		// 账户db初始化
		if config == nil {
			globalConfig, err := conf.GlobalConfig()
			if err != nil {
				log.Fatal(err)
			}
			config = globalConfig.Database
		}
		dbInstance, err = gorm.Open(config.Type, config.ConnectionStr())
		if err != nil {
			return
		}
		if config.LogMode {
			dbInstance.LogMode(true)
		}
	})
	return err
}

// 根据name获取db实例
func DBInstance() *gorm.DB {
	if dbInstance == nil {
		err := Open(nil)
		if err != nil {
			log.Fatal(err)
		}
	}
	return dbInstance
}

func switchDB(tx ...*gorm.DB) (*gorm.DB, error) {
	var db *gorm.DB
	switch len(tx) {
	case 0:
		db = DBInstance()
	default:
		db = tx[0]
	}
	if db == nil {
		return nil, response.ErrEmptyDB
	}
	return db, nil
}

func Begin() *gorm.DB {
	return DBInstance().Begin()
}
