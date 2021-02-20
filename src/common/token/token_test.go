package token

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"myProject/conf"
)

func TestAddToken(t *testing.T) {
	conf.InitConf("")
	tk, err := AddToken(context.Background(), uint64(10086))
	assert.Nil(t, err)
	t.Logf("token is %v", tk)
	assert.True(t, len(tk) > 0)
}
