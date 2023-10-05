package log

import (
	"context"
	"fmt"
	"github.com/golibs-starter/golib/log/field"
	"os"
	"strconv"
	"sync"
)

var _global *ZapLogger
var _globalLoggerLock = &sync.RWMutex{}

func init() {
	// We want the default global Logger will
	// in the same config with default in the log.Properties
	devMode, _ := strconv.ParseBool(os.Getenv("APP_LOGGING_DEVELOPMENT"))
	jsonMode := true
	jsonModeStr := os.Getenv("APP_LOGGING_JSONOUTPUTMODE")
	if jsonModeStr != "" {
		jsonMode, _ = strconv.ParseBool(jsonModeStr)
	}
	zapLogger, err := NewZapLogger(&Options{
		CallerSkip:     1,
		Development:    devMode,
		JsonOutputMode: jsonMode,
	})
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

func WithError(err error) Logger {
	return _global.Clone(-1).WithError(err)
}

func WithErrors(errs ...error) Logger {
	return _global.Clone(-1).WithErrors(errs...)
}

func WithAny(key string, value interface{}) Logger {
	return _global.Clone(-1).WithAny(key, value)
}

// Debug uses fmt.Sprint to construct and log a message.
func Debug(args ...interface{}) {
	_global.Debug(args...)
}

// Debugf uses fmt.Sprintf to log a template message.
func Debugf(format string, args ...interface{}) {
	_global.Debugf(format, args...)
}

// Debugc use WithCtx and fmt.Sprintf to log a template message.
func Debugc(ctx context.Context, msgFormat string, args ...interface{}) {
	_global.Debugc(ctx, msgFormat, args...)
}

// Info uses fmt.Sprint to construct and log a message.
func Info(args ...interface{}) {
	_global.Info(args...)
}

// Infof uses fmt.Sprintf to log a template message.
func Infof(msgFormat string, args ...interface{}) {
	_global.Infof(msgFormat, args...)
}

// Infoc use WithCtx and fmt.Sprintf to log a template message.
func Infoc(ctx context.Context, msgFormat string, args ...interface{}) {
	_global.Infoc(ctx, msgFormat, args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func Warn(args ...interface{}) {
	_global.Warn(args...)
}

// Warnf uses fmt.Sprintf to log a template message.
func Warnf(msgFormat string, args ...interface{}) {
	_global.Warnf(msgFormat, args...)
}

// Warnc use WithCtx and fmt.Sprintf to log a template message.
func Warnc(ctx context.Context, msgFormat string, args ...interface{}) {
	_global.Warnc(ctx, msgFormat, args...)
}

// Error uses fmt.Sprint to construct and log a message.
func Error(args ...interface{}) {
	_global.Error(args...)
}

// Errorf uses fmt.Sprintf to log a template message.
func Errorf(msgFormat string, args ...interface{}) {
	_global.Errorf(msgFormat, args...)
}

// Errorc use WithCtx and fmt.Sprintf to log a template message.
func Errorc(ctx context.Context, msgFormat string, args ...interface{}) {
	_global.Errorc(ctx, msgFormat, args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(args ...interface{}) {
	_global.Fatal(args...)
}

// Fatalf uses fmt.Sprintf to log a template message, then calls os.Exit.
func Fatalf(msgFormat string, args ...interface{}) {
	_global.Fatalf(msgFormat, args...)
}

// Fatalc use WithCtx and fmt.Sprintf to log a template message.
func Fatalc(ctx context.Context, msgFormat string, args ...interface{}) {
	_global.Fatalc(ctx, msgFormat, args...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func Panic(args ...interface{}) {
	_global.Panic(args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func Panicf(msgFormat string, args ...interface{}) {
	_global.Panicf(msgFormat, args...)
}

// Panicc use WithCtx and fmt.Sprintf to log a template message.
func Panicc(ctx context.Context, msgFormat string, args ...interface{}) {
	_global.Panicc(ctx, msgFormat, args...)
}
