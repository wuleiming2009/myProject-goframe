package log

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/coreos/go-systemd/v22/journal"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logging loggingZ
var once sync.Once

const (
	debugLevel uint = iota + 1
	infoLevel
	warnLevel
	errorLevel
	fatalLevel
	alterLevel
)

const (
	w byte = 'w'
	f byte = 'f'
)

func GetDebugLevel() uint {
	return debugLevel
}

func GetInfoLevel() uint {
	return infoLevel
}

func GetWarnLevel() uint {
	return warnLevel
}

func GetErrorLevel() uint {
	return errorLevel
}

func GetFatalLevel() uint {
	return fatalLevel
}

func GetAlertLevel() uint {
	return alterLevel
}

func GetFormatType() byte {
	return f
}

func GetWriteType() byte {
	return w
}

//var at atomicLevelManager

// basic logging
type loggingZ struct {
	consoleLogger  *zap.Logger
	standardLogger *zap.Logger
	atomicLevel    bool
}

func init() {
	flag.StringVar(&zlogDir, "z", "", "zap log dir")
	//flag.BoolVar(&logging.atomicLevel, "al", true, "zap atomic level")

	if journal.Enabled() {
		//nolint:errcheck
		zap.RegisterSink("journal", func(*url.URL) (zap.Sink, error) {
			return &journalSink{}, nil
		})
	}

	setConsoleLogger()
	//设置日志调节
	//at = NewAtomicLevelManager()
}

func setConsoleLogger() {
	// Mirrored from zap.NewProductionConfig
	cfg := &zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	if os.Getenv("LOG_TO_JOURNALD") != "" && journal.Enabled() {
		cfg.OutputPaths = []string{"journal://null"}
		cfg.ErrorOutputPaths = []string{"journal://null"}
	}

	logger, err := cfg.Build()
	if err != nil {
		fmt.Printf("new console loggger failed,err:%s \n", err)
		return
	}
	logging.consoleLogger = logger
}

func setStandardLogger() {
	logging.standardLogger = createLogger()
}

// start only once
func createLogger() *zap.Logger {
	debug := &syncBuffer{
		sev: 1,
	}
	//nolint:errcheck
	debug.rotateFile(time.Now())
	//定时将缓存中的日志刷入硬盘
	go debug.FlushDaemon()

	output := []zapcore.WriteSyncer{debug}
	output = append(output, os.Stdout)

	le := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	//if logging.atomicLevel {
	//	//开关打开再运行调节器
	//	at.run()
	//	le = at.GetLevel()
	//}

	enc := zap.NewProductionEncoderConfig()
	enc.TimeKey = "time"
	enc.CallerKey = "path_line"
	enc.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(zapcore.NewJSONEncoder(enc),
		zapcore.NewMultiWriteSyncer(output...),
		le,
	)

	var option []zap.Option
	option = append(option, zap.AddCaller())

	return zap.New(core, option...)
}

func (l *loggingZ) outputF(tag, depth uint, wt byte, msg string, args ...interface{}) {
	if !flag.Parsed() {
		l.consoleLogger.WithOptions(skip(depth+1)).Sugar().Infof(msg, args...)
		return
	}

	once.Do(func() {
		setStandardLogger()
	})

	switch tag {
	case fatalLevel:
		l.loggingFatal(depth+1, wt, msg, args...)
	case errorLevel:
		l.loggingError(depth+1, wt, msg, args...)
	case warnLevel:
		l.loggingWarn(depth+1, wt, msg, args...)
	case debugLevel:
		l.loggingDebug(depth+1, wt, msg, args...)
	case infoLevel:
		l.loggingInfo(depth+1, wt, msg, args...)
	default:
		l.loggingInfo(depth+1, wt, msg, args...)
	}
}

func (l *loggingZ) output(tag, depth uint, msg string, fields ...zap.Field) {
	if !flag.Parsed() {
		l.consoleLogger.WithOptions(skip(depth+1)).Info(msg, fields...)
		return
	}

	once.Do(func() {
		setStandardLogger()
	})

	switch tag {
	case fatalLevel:
		l.fatal(depth+1, msg, fields...)
	case errorLevel:
		l.error(depth+1, msg, fields...)
	case warnLevel:
		l.warn(depth+1, msg, fields...)
	case debugLevel:
		l.debug(depth+1, msg, fields...)
	case infoLevel:
		l.info(depth+1, msg, fields...)
	default:
		l.info(depth+1, msg, fields...)
	}
}

func (l *loggingZ) loggingDebug(depth uint, wt byte, msg string, args ...interface{}) {
	if wt == w {
		l.standardLogger.WithOptions(skip(depth+1)).Sugar().Debugw(msg, args...)
		return
	}
	l.standardLogger.WithOptions(skip(depth+1)).Sugar().Debugf(msg, args...)
}

func (l *loggingZ) loggingInfo(depth uint, wt byte, msg string, args ...interface{}) {
	if wt == w {
		l.standardLogger.WithOptions(skip(depth+1)).Sugar().Infow(msg, args...)
		return
	}
	l.standardLogger.WithOptions(skip(depth+1)).Sugar().Infof(msg, args...)
}

func (l *loggingZ) loggingWarn(depth uint, wt byte, msg string, args ...interface{}) {
	if wt == w {
		l.standardLogger.WithOptions(skip(depth+1)).Sugar().Warnw(msg, args...)
		return
	}
	l.standardLogger.WithOptions(skip(depth+1)).Sugar().Warnf(msg, args...)
}

func (l *loggingZ) loggingError(depth uint, wt byte, msg string, args ...interface{}) {
	if wt == w {
		l.standardLogger.WithOptions(skip(depth+1)).Sugar().Errorw(msg, args...)
		return
	}
	l.standardLogger.WithOptions(skip(depth+1)).Sugar().Errorf(msg, args...)
}

func (l *loggingZ) loggingFatal(depth uint, wt byte, msg string, args ...interface{}) {
	//nolint:errcheck
	l.standardLogger.Sync()
	if wt == w {
		l.standardLogger.WithOptions(skip(depth+1)).Sugar().Fatalw(msg, args...)
		return
	}
	l.standardLogger.WithOptions(skip(depth+1)).Sugar().Fatalf(msg, args...)
}

func (l *loggingZ) debug(depth uint, msg string, fields ...zap.Field) {
	l.standardLogger.WithOptions(skip(depth+1)).Debug(msg, fields...)
}

func (l *loggingZ) info(depth uint, msg string, fields ...zap.Field) {
	l.standardLogger.WithOptions(skip(depth+1)).Info(msg, fields...)
}

func (l *loggingZ) warn(depth uint, msg string, fields ...zap.Field) {
	l.standardLogger.WithOptions(skip(depth+1)).Warn(msg, fields...)
}

func (l *loggingZ) error(depth uint, msg string, fields ...zap.Field) {
	l.standardLogger.WithOptions(skip(depth+1)).Error(msg, fields...)
}

func (l *loggingZ) fatal(depth uint, msg string, fields ...zap.Field) {
	l.standardLogger.WithOptions(skip(depth+1)).Fatal(msg, fields...)
}

func skip(depth uint) zap.Option {
	return zap.AddCallerSkip(int(depth))
}

func Flush() {
	if logging.standardLogger != nil {
		//nolint:errcheck
		logging.standardLogger.Sync()
	}
}
