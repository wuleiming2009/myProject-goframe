package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"myProject/conf"
)

func TestFetchWithJson(t *testing.T) {
	ast := assert.New(t)
	conf.InitConf("")
	type testFetch struct {
		A int `json:"a"`
	}
	redisCli, err := RedisClient()
	ast.NoError(err)
	key := "test_fetch"
	tmp := &testFetch{}
	err = FetchWithJson(context.Background(), redisCli, key, 1*time.Minute, tmp, func() (interface{}, error) {
		return &testFetch{A: 1}, nil
	})
	ast.NoError(err)
	valueFromRedis := redisCli.Get(context.Background(), key).String()
	t.Logf("valueFromRedis, %s", valueFromRedis)
	ast.NotEmpty(valueFromRedis)
	ast.EqualValues(tmp, &testFetch{A: 1})
}
