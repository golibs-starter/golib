package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
)

type TestingLogger struct {
	*DefaultLogger
	tb testing.TB
}

func NewTestingLogger(tb testing.TB, options *Options) (*TestingLogger, error) {
	defaultLogger, err := NewDefaultLogger(options)
	if err != nil {
		return nil, err
	}
	return NewTestingLoggerFromDefault(tb, defaultLogger), nil
}

func NewTestingLoggerFromDefault(tb testing.TB, defaultLogger *DefaultLogger) *TestingLogger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	coreLogger := defaultLogger.coreLogger.WithOptions(
		zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderConfig),
				NewTestingWriter(tb),
				zap.NewAtomicLevelAt(zap.DebugLevel),
			)
		}),
	)
	return &TestingLogger{
		DefaultLogger: &DefaultLogger{
			options:     defaultLogger.options,
			coreLogger:  coreLogger,
			sugarLogger: coreLogger.Sugar(),
		},
		tb: tb,
	}
}
