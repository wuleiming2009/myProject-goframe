package models

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/jmoiron/sqlx"

	"myProject/conf"
)

const MYSQL_INIT = "root:root@tcp(mysql:3306)/"
const MysqlInit = "root:root@tcp(mysql:3306)/unit_test"

func conDB() (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", MysqlInit)
	if err != nil {
		fmt.Println("DB initialize failed", err)
		return nil, err
	}
	return db, err
}

func runSqlFromFiles(fileNames []string) error {
	db, err := conDB()
	if err != nil {
		return err
	}
	_, _ = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName + " DEFAULT CHARACTER SET utf8mb4;")
	_, _ = db.Exec("USE " + dbName + ";")

	for _, fileName := range fileNames {
		file, ioErr := ioutil.ReadFile(fileName)
		if ioErr != nil {
			fmt.Println("Open db file error", fileName, ioErr)
			return ioErr
		}
		requests := strings.Split(string(file), ";\n")
		for _, request := range requests {
			rawRequest := strings.TrimSuffix(request, "\n")
			if rawRequest == "" {
				continue
			}
			_, err := db.Exec(request)
			if err != nil {
				fmt.Println("Sql execution failed", err)
				return err
			}
		}
	}
	return nil

}

func initDBMode() {
	db, err := conDB()
	if err != nil {
		log.Println(err)
		return
	}
	_, _ = db.Exec("SET global sql_mode = '';")
	defer db.Close()
}

func initDBInfo(dbfile string) {
	initDBMode()
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	fileName := filepath.Join("../migration", dbfile)
	_ = runSqlFromFiles([]string{fileName})
	time.Sleep(1 * time.Second)
}

func ClearTestDB() error {
	db, err := conDB()
	if err != nil {
		return err
	}
	_, err = db.Exec("DROP DATABASE `" + dbName + "`;")
	log.Println("clear DB.")
	return err
}

var dbName = "dress"
var dbConn = MYSQL_INIT + dbName
var db *sqlx.DB
var dbConfig = &conf.DBConfig{
	Type:     "mysql",
	User:     "root",
	Password: "root",
	Host:     "mysql",
	Db:       dbName,
	LogMode:  true,
}

func InitTestDB() {
	initDBInfo("tb_goods.sql")
	time.Sleep(1 * time.Second)
	var err error
	err = Open(dbConfig)
	if err != nil {
		glog.Fatal(err)
	}
}
