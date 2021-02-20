package conf

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/go-ini/ini"

	"myProject/lib/log"
)

type ServerConfig struct {
	HttpPort     int32         `json:"http_port" ini:"http_port"`
	ReadTimeout  time.Duration `json:"read_timeout" ini:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout" ini:"write_timeout"`
}

type DBConfig struct {
	Type     string `json:"type" ini:"type"`
	User     string `json:"user" ini:"user"`
	Password string `json:"password" ini:"password"`
	Host     string `json:"host" ini:"host"`
	Db       string `json:"db" ini:"db"`
	LogMode  bool   `json:"log_mode" ini:"log_mode"`
	Loc      string `json:"loc" ini:"loc"`
}

func (c *DBConfig) ConnectionStr() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		c.User,
		c.Password,
		c.Host,
		c.Db,
		c.Loc)
}

type RedisConfig struct {
	Addr     string `json:"addr" ini:"addr"`
	Password string `json:"password" ini:"password"`
	DB       int    `json:"db" ini:"db"`
}
type AppConfig struct {
	JwtSecret       string        `json:"jwt_secret" ini:"jwt_secret"`
	TokenExpireTime time.Duration `json:"token_expire_time" ini:"token_expire_time"`
}
type Config struct {
	RunMode  string        `json:"run_mode" ini:"run_mode"`
	Server   *ServerConfig `json:"server" ini:"server"`
	Database *DBConfig     `json:"database" ini:"database"`
	Redis    *RedisConfig  `json:"redis" ini:"redis"`
	App      *AppConfig    `json:"app" ini:"app"`
}

type PaypalConfig struct {
	ClientId string `json:"client_id" ini:"client_id"`
	SecretId string `json:"secret_id" ini:"secret_id"`
	ApiBase  string `json:"api_base" ini:"api_base"`
}

var globalConfig *Config
var once sync.Once

func InitConf(cfgPath string) {
	if cfgPath == "" {
		_, dir, _, _ := runtime.Caller(0)
		cfgPath = filepath.Join(filepath.Dir(dir), "app.ini")
	}
	once.Do(func() {
		cfg, err := ini.Load(cfgPath)
		if err != nil {
			log.Fatalf("Fail to parse '%s': %v", cfgPath, err)
		}
		globalConfig = &Config{}
		err = cfg.MapTo(globalConfig)
		if err != nil {
			log.Fatal(err)
		}
	})
}

func GlobalConfig() (*Config, error) {
	if globalConfig == nil {
		return nil, errors.New("global config is not initialized")
	}
	return globalConfig, nil
}

func GetAppConfig() (*AppConfig, error) {
	cfg, err := GlobalConfig()
	if err != nil {
		return nil, err
	}
	if cfg.App == nil {
		return nil, errors.New("empty app config")
	}
	return cfg.App, nil
}
