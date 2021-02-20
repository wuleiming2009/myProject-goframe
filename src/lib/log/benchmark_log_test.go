//go test -bench=. -benchmem
//goos: darwin
//goarch: amd64
//pkg: ptapp.cn/util/zaplog
//Benchmark_zaplog-4      10000000               451 ns/op               4 B/op          0 allocs/op
//Benchmark_glog-4         1000000              5114 ns/op             600 B/op          9 allocs/op
package log

import (
	"fmt"
	"strings"
	"testing"

	"github.com/golang/glog"
	"go.uber.org/zap"
)

func fakeMessages(n int) string {
	messages := make([]string, n)
	for i := range messages {
		messages[i] = fmt.Sprintf("Test logging, but use a somewhat realistic message length. (#%v)", i)
	}
	return strings.Join(messages, "")
}

func Benchmark_glog(b *testing.B) {
	str := fakeMessages(1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		glog.Info(str)
		glog.Info(str)
		glog.Warning(str)
	}
	glog.Flush()
}

func Benchmark_zaplog(b *testing.B) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{}

	log, _ := cfg.Build()
	str := fakeMessages(1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info(str)
		log.Debug(str)
		log.Warn(str)
	}
	//nolint:errcheck
	log.Sync()
}
