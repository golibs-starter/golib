package log

import (
	"context"
	"fmt"
)

var global Logger

func init() {
	var err error
	if global, err = NewZapLogger(&Options{CallerSkip: 2}); err != nil {
		panic(fmt.Errorf("init global logger error [%v]", err))
	}
}

// ReplaceGlobal Register a logger instance as global
func ReplaceGlobal(logger Logger) {
	global = logger
}

// GetGlobal Get global logger instance
func GetGlobal() Logger {
	return global
}

func WithCtx(ctx context.Context) Logger {
	return global.WithCtx(ctx)
}

func WithField(fields ...Field) Logger {
	return global.WithField(fields...)
}

// Info uses fmt.Sprint to construct and log a message.
func Info(args ...interface{}) {
	global.Info(args...)
}

// Infof uses fmt.Sprintf to log a template message.
func Infof(msgFormat string, args ...interface{}) {
	global.Infof(msgFormat, args...)
}

// Debug uses fmt.Sprint to construct and log a message.
func Debug(args ...interface{}) {
	global.Debug(args...)
}

// Debugf uses fmt.Sprintf to log a template message.
func Debugf(msgFormat string, args ...interface{}) {
	global.Debugf(msgFormat, args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func Warn(args ...interface{}) {
	global.Warn(args...)
}

// Warnf uses fmt.Sprintf to log a template message.
func Warnf(msgFormat string, args ...interface{}) {
	global.Warnf(msgFormat, args...)
}

// Error uses fmt.Sprint to construct and log a message.
func Error(args ...interface{}) {
	global.Error(args...)
}

// Errorf uses fmt.Sprintf to log a template message.
func Errorf(msgFormat string, args ...interface{}) {
	global.Errorf(msgFormat, args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(args ...interface{}) {
	global.Fatal(args...)
}

// Fatalf uses fmt.Sprintf to log a template message, then calls os.Exit.
func Fatalf(msgFormat string, args ...interface{}) {
	global.Fatalf(msgFormat, args...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func Panic(args ...interface{}) {
	global.Panic(args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func Panicf(msgFormat string, args ...interface{}) {
	global.Panicf(msgFormat, args...)
}
