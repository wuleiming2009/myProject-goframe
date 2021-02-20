package log

import (
	"context"
	"flag"

	"go.uber.org/zap"
)

func Info(msg string, fields ...zap.Field) {
	logging.output(infoLevel, 1, msg, fields...)
}

func Infof(msg string, args ...interface{}) {
	logging.outputF(infoLevel, 1, f, msg, args...)
}

func Infow(msg string, kvPairs ...interface{}) {
	logging.outputF(infoLevel, 1, w, msg, kvPairs...)
}

func Debug(msg string, fields ...zap.Field) {
	logging.output(debugLevel, 1, msg, fields...)
}

func Debugf(msg string, args ...interface{}) {
	logging.outputF(debugLevel, 1, f, msg, args...)
}

func Debugw(msg string, kvPairs ...interface{}) {
	logging.outputF(debugLevel, 1, w, msg, kvPairs...)
}

func Warn(msg string, fields ...zap.Field) {
	logging.output(warnLevel, 1, msg, fields...)
}

func Warnf(msg string, args ...interface{}) {
	logging.outputF(warnLevel, 1, f, msg, args...)
}

func Warnw(msg string, kvPairs ...interface{}) {
	logging.outputF(warnLevel, 1, w, msg, kvPairs...)
}

func Error(msg string, fields ...zap.Field) {
	logging.output(errorLevel, 1, msg, fields...)
}

func Errorf(msg string, args ...interface{}) {
	logging.outputF(errorLevel, 1, f, msg, args...)
}

func Errorw(msg string, kvPairs ...interface{}) {
	logging.outputF(errorLevel, 1, w, msg, kvPairs...)
}

func Alert(msg string, fields ...zap.Field) {
	logging.output(alterLevel, 1, msg, fields...)
}

func Alertf(msg string, args ...interface{}) {
	logging.outputF(alterLevel, 1, f, msg, args...)
}

func Alertw(msg string, kvPairs ...interface{}) {
	logging.outputF(alterLevel, 1, w, msg, kvPairs...)
}

func Fatal(args ...interface{}) {
	logging.outputF(fatalLevel, 1, f, "fatal message : %v", args...)
}

func Fatalf(msg string, args ...interface{}) {
	logging.outputF(fatalLevel, 1, f, msg, args...)
}

func Fatalw(msg string, kvPairs ...interface{}) {
	logging.outputF(fatalLevel, 1, w, msg, kvPairs...)
}

type ZapLogger struct {
	Fields []zap.Field
	z      *loggingZ
}

func NewZapLogger(reqid string) *ZapLogger {
	var fields []zap.Field
	fields = append(fields, ReqId(reqid), MachineField(host))
	return &ZapLogger{
		Fields: fields,
		z: &loggingZ{
			consoleLogger: logging.consoleLogger,
		},
	}
}

func NewLogger(reqid string) *ZapLogger {
	return NewZapLogger(reqid)
}

func NewContext(ctx context.Context, reqId string) context.Context {
	if ctx == nil {
		ctx = context.TODO()
	}

	zapLog := NewZapLogger(reqId)
	return context.WithValue(ctx, GetZapKey(), zapLog)
}

func FromContext(ctx context.Context) *ZapLogger {
	if ctx == nil {
		return NewZapLogger("")
	}

	zl, ok := ctx.Value(GetZapKey()).(*ZapLogger)
	if !ok {
		zl = NewZapLogger("")
		ctx = context.WithValue(ctx, GetZapKey(), zl)
	}
	return zl
}

func (t *ZapLogger) Output(tag, depth uint, msg string, field ...zap.Field) {
	if t.z.standardLogger != nil {
		t.z.output(tag, depth, msg, field...)
		return
	}

	if flag.Parsed() {
		if logging.standardLogger == nil {
			once.Do(func() {
				setStandardLogger()
			})
		}

		t.z.standardLogger = logging.standardLogger.With(t.Fields...)
		t.z.output(tag, depth, msg, field...)
		return
	}

	t.z.output(tag, depth, msg, field...)
}

func (t *ZapLogger) OutputF(tag, depth uint, wt byte, msg string, args ...interface{}) {
	if t.z.standardLogger != nil {
		t.z.outputF(tag, depth, wt, msg, args...)
		return
	}

	if flag.Parsed() {
		if logging.standardLogger == nil {
			once.Do(func() {
				setStandardLogger()
			})
		}

		t.z.standardLogger = logging.standardLogger.With(t.Fields...)
		t.z.outputF(tag, depth, wt, msg, args...)
		return
	}

	t.z.outputF(tag, depth, wt, msg, args...)
}

func (t *ZapLogger) Debug(msg string, fields ...zap.Field) {
	t.Output(debugLevel, 2, msg, fields...)
}

func (t *ZapLogger) Debugf(msg string, args ...interface{}) {
	t.OutputF(debugLevel, 2, f, msg, args...)
}

func (t *ZapLogger) Debugw(msg string, kvPairs ...interface{}) {
	t.OutputF(debugLevel, 2, w, msg, kvPairs...)
}

func (t *ZapLogger) Info(msg string, fields ...zap.Field) {
	t.Output(infoLevel, 2, msg, fields...)
}

func (t *ZapLogger) Infof(msg string, args ...interface{}) {
	t.OutputF(infoLevel, 2, f, msg, args...)
}

func (t *ZapLogger) Infow(msg string, kvPairs ...interface{}) {
	t.OutputF(infoLevel, 2, w, msg, kvPairs...)
}

func (t *ZapLogger) Warn(msg string, fields ...zap.Field) {
	t.Output(warnLevel, 2, msg, fields...)
}

func (t *ZapLogger) Warnf(msg string, args ...interface{}) {
	t.OutputF(warnLevel, 2, f, msg, args...)
}

func (t *ZapLogger) Warnw(msg string, kvPairs ...interface{}) {
	t.OutputF(warnLevel, 2, w, msg, kvPairs...)
}

func (t *ZapLogger) Error(msg string, fields ...zap.Field) {
	t.Output(errorLevel, 2, msg, fields...)
}

func (t *ZapLogger) Errorf(msg string, args ...interface{}) {
	t.OutputF(errorLevel, 2, f, msg, args...)
}

func (t *ZapLogger) Errorw(msg string, kvPairs ...interface{}) {
	t.OutputF(errorLevel, 2, w, msg, kvPairs...)
}

func (t *ZapLogger) Fatal(msg string, fields ...zap.Field) {
	t.Output(errorLevel, 2, msg, fields...)
}

func (t *ZapLogger) Fatalf(msg string, args ...interface{}) {
	t.OutputF(fatalLevel, 2, f, msg, args...)
}

func (t *ZapLogger) Fatalw(msg string, kvPairs ...interface{}) {
	t.OutputF(fatalLevel, 2, w, msg, kvPairs...)
}

func (t *ZapLogger) Flush() {
	if t.z.standardLogger != nil {
		//nolint:errcheck
		t.z.standardLogger.Sync()
	}
}

type nopLogger struct{}

// NewNopLogger returns a logger that doesn't do anything.
func NewNopLogger() nopLogger { return nopLogger{} }

func (nopLogger) Log(...interface{}) error { return nil }
