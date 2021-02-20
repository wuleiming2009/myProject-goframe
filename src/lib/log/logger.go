package log

import (
	"fmt"

	"github.com/golang/glog"
)

type GLogger struct {
	reqid  string
	format string
}

var (
	alertPrefix = "[ MT ALERT ] "
	errorPrefix = "[ MT ERROR ] "
)

func New(reqid string) *GLogger {
	l := &GLogger{reqid: reqid}
	l.format = fmt.Sprintf("[%s] ", l.reqid)
	return l
}

func (l *GLogger) Infof(format string, args ...interface{}) {
	glog.InfoDepth(1, fmt.Sprintf(format, args...))
}

func (l *GLogger) InfofDepth(depth int, format string, args ...interface{}) {
	glog.InfoDepth(depth+1, l.format, fmt.Sprintf(format, args...))
}

func (l *GLogger) Info(args ...interface{}) {
	glog.InfoDepth(1, l.format, fmt.Sprint(args...))
}

func (l *GLogger) InfoDepth(depth int, args ...interface{}) {
	glog.InfoDepth(depth+1, l.format, fmt.Sprint(args...))
}

func (l *GLogger) Warnf(format string, args ...interface{}) {
	glog.WarningDepth(1, fmt.Sprintf(format, args...))
}

func (l *GLogger) WarnfDepth(depth int, format string, args ...interface{}) {
	glog.WarningDepth(depth+1, l.format, fmt.Sprintf(format, args...))
}

func (l *GLogger) WarnDepth(depth int, args ...interface{}) {
	glog.WarningDepth(depth+1, l.format, fmt.Sprint(args...))
}

func (l *GLogger) Warn(args ...interface{}) {
	glog.WarningDepth(1, l.format, fmt.Sprint(args...))
}

func (l *GLogger) Alertf(format string, args ...interface{}) {
	glog.ErrorDepth(1, fmt.Sprintf(alertPrefix+format, args...))
}

func (l *GLogger) Alert(args ...interface{}) {
	glog.ErrorDepth(1, alertPrefix, fmt.Sprint(args...))
}

func (l *GLogger) Errorf(format string, args ...interface{}) {
	glog.ErrorDepth(1, errorPrefix, fmt.Sprintf(format, args...))
}

func (l *GLogger) ErrorfDepth(depth int, format string, args ...interface{}) {
	glog.ErrorDepth(depth+1, errorPrefix+l.format, fmt.Sprintf(format, args...))
}

func (l *GLogger) ErrorDepth(depth int, args ...interface{}) {
	glog.ErrorDepth(depth+1, errorPrefix+l.format, fmt.Sprint(args...))
}

func (l *GLogger) Error(args ...interface{}) {
	glog.ErrorDepth(1, errorPrefix+l.format, fmt.Sprint(args...))
}

func (l *GLogger) FatalDepth(depth int, args ...interface{}) {
	glog.FatalDepth(depth+1, errorPrefix+l.format, fmt.Sprint(args...))
}

func (l *GLogger) Fatal(args ...interface{}) {
	glog.FatalDepth(1, l.format, fmt.Sprint(args...))
}

func (l *GLogger) GetReqId() string {
	return l.reqid
}

func (l *GLogger) GetFormat() string {
	return l.format
}

func (l *GLogger) Verbose(level int32) Verbose {
	return Verbose{
		GLogger: l,
		v:       glog.V(glog.Level(level)),
	}
}

type Verbose struct {
	*GLogger
	v glog.Verbose
}

func (v Verbose) V() glog.Verbose {
	return v.v
}
