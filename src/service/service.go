package service

import (
	"sync"

	"myProject/lib/log"
)

var once sync.Once

var (
	serviceAccount AccountService
)

// 初始化service单例
func init() {
	var err error
	serviceAccount, err = NewAccountService()
	if err != nil {
		log.Fatal(err)
	}
}
