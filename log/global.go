package log

import (
	"context"
	"fmt"
	"gitlab.com/golibs-starter/golib/log/field"
	"sync"
)

var _global *ZapLogger
var _globalLoggerLock = &sync.RWMutex{}

func init() {
	zapLogger, err := NewZapLogger(&Options{CallerSkip: 1, Development: true})
	if err != nil {
		panic(fmt.Errorf("init global logger error [%v]", err))
	}
	ReplaceGlobal(zapLogger)
}

// ReplaceGlobal Register a logger instance as global
func ReplaceGlobal(logger *ZapLogger) {
	_globalLoggerLock.Lock()
	defer _globalLoggerLock.Unlock()
	_global = logger.Clone(1)
}

// GetGlobal Get global logger instance
func GetGlobal() Logger {
	return _global
}

func WithCtx(ctx context.Context, additionalFields ...field.Field) Logger {
	return _global.Clone(-1).WithCtx(ctx, additionalFields...)
}

func WithField(fields ...field.Field) Logger {
	return _global.Clone(-1).WithField(fields...)
}

func WithErrors(errs ...error) Logger {
	return _global.Clone(-1).WithErrors(errs...)
}

func WithAny(key string, value interface{}) Logger {
	return _global.Clone(-1).WithAny(key, value)
}

// Info uses fmt.Sprint to construct and log a message.
func Info(args ...interface{}) {
	_global.Info(args...)
}

// Infof uses fmt.Sprintf to log a template message.
func Infof(msgFormat string, args ...interface{}) {
	_global.Infof(msgFormat, args...)
}

// Debug uses fmt.Sprint to construct and log a message.
func Debug(args ...interface{}) {
	_global.Debug(args...)
}

// Debugf uses fmt.Sprintf to log a template message.
func Debugf(msgFormat string, args ...interface{}) {
	_global.Debugf(msgFormat, args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func Warn(args ...interface{}) {
	_global.Warn(args...)
}

// Warnf uses fmt.Sprintf to log a template message.
func Warnf(msgFormat string, args ...interface{}) {
	_global.Warnf(msgFormat, args...)
}

// Error uses fmt.Sprint to construct and log a message.
func Error(args ...interface{}) {
	_global.Error(args...)
}

// Errorf uses fmt.Sprintf to log a template message.
func Errorf(msgFormat string, args ...interface{}) {
	_global.Errorf(msgFormat, args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(args ...interface{}) {
	_global.Fatal(args...)
}

// Fatalf uses fmt.Sprintf to log a template message, then calls os.Exit.
func Fatalf(msgFormat string, args ...interface{}) {
	_global.Fatalf(msgFormat, args...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func Panic(args ...interface{}) {
	_global.Panic(args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func Panicf(msgFormat string, args ...interface{}) {
	_global.Panicf(msgFormat, args...)
}
