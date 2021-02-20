package cache

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"

	"myProject/conf"
	"myProject/lib/log"
)

var once sync.Once
var redisCli *redis.Client

func NewRedisClient(config *conf.RedisConfig) {
	var err error
	once.Do(func() {
		if config == nil {
			cfg, err := conf.GlobalConfig()
			if err != nil {
				log.Fatal(err)
			}
			config = cfg.Redis
		}
		log.Infof("Redis Config:%+v", config)
		redisCli = redis.NewClient(&redis.Options{
			Addr:     config.Addr,
			Password: config.Password,
			DB:       config.DB,
		})
		err = redisCli.Ping(context.Background()).Err()
		if err != nil {
			log.Fatal(err)
			return
		}
	})
}

func RedisClient() (*redis.Client, error) {
	if redisCli == nil {
		NewRedisClient(nil)
	}
	if redisCli == nil {
		return nil, errors.New("empty redisCli")
	}
	return redisCli, nil
}

func CloseRedis() error {
	if redisCli != nil {
		return redisCli.Close()
	}
	return nil
}

func FetchWithJson(ctx context.Context, redisCli *redis.Client, key string, expiration time.Duration, dest interface{}, fn func() (interface{}, error)) error {
	zl := log.FromContext(ctx)
	value, err := redisCli.Get(ctx, key).Result()
	if err == redis.Nil {
		ret, err := fn()
		if err != nil {
			zl.Errorf("fn execute failed, %v", err)
			return err
		}
		marshalValue, err := json.Marshal(ret)
		if err != nil {
			zl.Errorf("Marshal ret failed, %+v", err)
			return err
		}
		err = redisCli.Set(ctx, key, marshalValue, expiration).Err()
		if err != nil {
			zl.Errorf("save to redis failed, key:%s, marshalValue:%s, %v", key, marshalValue, err)
			return err
		}
		value = string(marshalValue)
	} else if err != nil {
		zl.Errorf("get from redis failed, key:%s, %v", key, err)
		return err
	}
	err = json.Unmarshal([]byte(value), dest)
	if err != nil {
		zl.Errorf("unmarshal value to ret failed, value:%s, %v", value, err)
		return err
	}
	return nil
}

func CleanCache(ctx context.Context, redisCli *redis.Client, key string) (err error) {
	if redisCli == nil {
		redisCli, err = RedisClient()
		if err != nil {
			return err
		}
	}
	return redisCli.Del(ctx, key).Err()
}
