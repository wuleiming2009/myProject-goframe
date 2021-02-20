package snow_flake

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"myProject/conf"
)

func Test_genSnowFlake(t *testing.T) {
	conf.InitConf("")
	ret := GenSnowFlake()
	t.Logf("genSnowFlake ret:%v", ret)
	assert.True(t, ret > 0)
}
