package log

import (
	"context"
	"fmt"
	"github.com/golibs-starter/golib/log"
	"github.com/golibs-starter/golib/pubsub"
	"sync"
)

var _global log.Logger
var _globalLoggerLock = &sync.RWMutex{}

func init() {
	var err error
	if _global, err = log.NewZapLogger(&log.Options{CallerSkip: 2}); err != nil {
		panic(fmt.Errorf("init global web logger error [%v]", err))
	}
}

// ReplaceGlobal Register a logger instance as global
func ReplaceGlobal(logger log.Logger) {
	_globalLoggerLock.Lock()
	defer _globalLoggerLock.Unlock()
	_global = logger
}

// Debug
// Deprecated: use `github.com/golibs-starter/golib/log.WithCtx(ctx).Debugf` instead
func Debug(ctx context.Context, msgFormat string, args ...interface{}) {
	_global.WithCtx(ctx).Debugf(msgFormat, args...)
}

// Info
// Deprecated: use `github.com/golibs-starter/golib/log.WithCtx(ctx).Infof` instead
func Info(ctx context.Context, msgFormat string, args ...interface{}) {
	_global.WithCtx(ctx).Infof(msgFormat, args...)
}

// Warn
// Deprecated: use `github.com/golibs-starter/golib/log.WithCtx(ctx).Warnf` instead
func Warn(ctx context.Context, msgFormat string, args ...interface{}) {
	_global.WithCtx(ctx).Warnf(msgFormat, args...)
}

// Error
// Deprecated: use `github.com/golibs-starter/golib/log.WithCtx(ctx).Errorf` instead
func Error(ctx context.Context, msgFormat string, args ...interface{}) {
	_global.WithCtx(ctx).Errorf(msgFormat, args...)
}

// Fatal
// Deprecated: use `github.com/golibs-starter/golib/log.WithCtx(ctx).Fatalf` instead
func Fatal(ctx context.Context, msgFormat string, args ...interface{}) {
	_global.WithCtx(ctx).Fatalf(msgFormat, args...)
}

// Debuge
// Deprecated: use `github.com/golibs-starter/golib/log.WithCtx(ctx).Debuge` instead
func Debuge(e pubsub.Event, msgFormat string, args ...interface{}) {
	_global.WithCtx(e.Context()).Debugf(msgFormat, args...)
}

// Infoe
// Deprecated: use `github.com/golibs-starter/golib/log.WithCtx(ctx).Infoe` instead
func Infoe(e pubsub.Event, msgFormat string, args ...interface{}) {
	_global.WithCtx(e.Context()).Infof(msgFormat, args...)
}

// Warne
// Deprecated: use `github.com/golibs-starter/golib/log.WithCtx(ctx).Warne` instead
func Warne(e pubsub.Event, msgFormat string, args ...interface{}) {
	_global.WithCtx(e.Context()).Warnf(msgFormat, args...)
}

// Errore
// Deprecated: use `github.com/golibs-starter/golib/log.WithCtx(ctx).Errore` instead
func Errore(e pubsub.Event, msgFormat string, args ...interface{}) {
	_global.WithCtx(e.Context()).Errorf(msgFormat, args...)
}

// Fatale
// Deprecated: use `github.com/golibs-starter/golib/log.WithCtx(ctx).Fatale` instead
func Fatale(e pubsub.Event, msgFormat string, args ...interface{}) {
	_global.WithCtx(e.Context()).Fatalf(msgFormat, args...)
}

// Debugf
// Deprecated: use `github.com/golibs-starter/golib/log.WithCtx(ctx).Debugf` instead
func Debugf(msgFormat string, args ...interface{}) {
	_global.Debugf(msgFormat, args...)
}

// Infof
// Deprecated: use `github.com/golibs-starter/golib/log.WithCtx(ctx).Infof` instead
func Infof(msgFormat string, args ...interface{}) {
	_global.Infof(msgFormat, args...)
}

// Warnf
// Deprecated: use `github.com/golibs-starter/golib/log.WithCtx(ctx).Warnf` instead
func Warnf(msgFormat string, args ...interface{}) {
	_global.Warnf(msgFormat, args...)
}

// Errorf
// Deprecated: use `github.com/golibs-starter/golib/log.WithCtx(ctx).Errorf` instead
func Errorf(msgFormat string, args ...interface{}) {
	_global.Errorf(msgFormat, args...)
}

// Fatalf
// Deprecated: use `github.com/golibs-starter/golib/log.WithCtx(ctx).Fatalf` instead
func Fatalf(msgFormat string, args ...interface{}) {
	_global.Fatalf(msgFormat, args...)
}
