package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestZapLog_concurrent(t *testing.T) {
	logger := zap.NewExample()
	logger.Info("hello my girl")
	nl := logger.WithOptions(zap.Fields(zap.String("reqid", "1234")))
	nl.Debug("hello my boy")
}

func TestNewZapLogger(t *testing.T) {
	ctx := NewContext(context.Background(), "testDl")
	zl := FromContext(ctx)
	zl.Errorf("test from context")
}

func TestNewContext(t *testing.T) {
	ctx := context.Background()

	zl := FromContext(ctx)
	assert.NotNil(t, zl)
	zl.Infof("hello")

	ReqCtx := NewContext(context.Background(), "testZap")
	l := FromContext(ReqCtx)
	l.Infof("hello")
}
