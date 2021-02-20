package routers

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"myProject/common/metadata"
	"myProject/conf"
	"myProject/controllers"
	_ "myProject/docs"
)

func InitRouter(cfg *conf.Config) *gin.Engine {
	gin.SetMode(cfg.RunMode)
	r := gin.New()

	// Logging to a file.
	accessLogPath, err := filepath.Abs("log")
	if err != nil {
		log.Fatal(err)
	}

	if gin.Mode() != gin.DebugMode {
		f, _ := os.Create(filepath.Join(accessLogPath, "access.log"))
		gin.DefaultWriter = io.MultiWriter(f)
		gin.DefaultErrorWriter = io.MultiWriter(f)
	}
	r.Use(gin.Logger())

	r.Use(gin.Recovery())
	r.Use(metadata.LogWithReqId)
	r.Use(metadata.Cors())

	// swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 账户系统路由
	account := r.Group("/account")
	{
		account.POST("/sign_up", controllers.AddAccount)
		account.POST("/login", controllers.UserLogin)
	}

	return r
}
