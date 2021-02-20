package controllers

import (
	"myProject/lib/log"
	"myProject/service"
)

// service实例
var (
	serviceAccount service.AccountService
)

func init() {
	var err error
	serviceAccount, err = service.NewAccountService()
	if err != nil {
		log.Fatal(err)
	}
}
