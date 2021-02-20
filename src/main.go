package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/glog"

	"myProject/cache"
	"myProject/conf"
	"myProject/lib/log"
	"myProject/routers"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	var cfgPath string
	flag.StringVar(&cfgPath, "f", "", "config file path")
	flag.Parse()
	conf.InitConf(cfgPath)
	config, err := conf.GlobalConfig()
	if err != nil {
		log.Fatal(err)
	}

	// 初始化路由
	router := routers.InitRouter(config)

	serverConfig := config.Server
	if config.Server == nil {
		log.Fatal("empty server config")
	}
	// 初始化http服务
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", serverConfig.HttpPort),
		Handler:        router,
		ReadTimeout:    time.Second * serverConfig.ReadTimeout,
		WriteTimeout:   time.Second * serverConfig.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	// 启动监听
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Info("Server exiting")
	beforeExit()
	log.Info("exit")
}

func beforeExit() {
	_ = cache.CloseRedis()
	glog.Flush()
	log.Flush()
}
