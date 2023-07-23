package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
)

type TestingLogger struct {
	*ZapLogger
	tb testing.TB
}

func NewTestingLogger(tb testing.TB, options *Options) (*TestingLogger, error) {
	defaultLogger, err := NewZapLogger(options)
	if err != nil {
		return nil, err
	}
	return NewTestingLoggerFromDefault(tb, defaultLogger), nil
}

func NewTestingLoggerFromDefault(tb testing.TB, defaultLogger *ZapLogger) *TestingLogger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	coreLogger := defaultLogger.core.WithOptions(
		zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderConfig),
				NewTestingWriter(tb),
				zap.NewAtomicLevelAt(zap.DebugLevel),
			)
		}),
	)
	return &TestingLogger{
		ZapLogger: &ZapLogger{
			core:  coreLogger,
			sugar: coreLogger.Sugar(),
		},
		tb: tb,
	}
}
