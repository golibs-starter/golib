package log

import (
	"context"
	"gitlab.com/golibs-starter/golib/log/field"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	opts  *Options
	core  *zap.Logger
	sugar *zap.SugaredLogger
}

func NewZapLogger(opts *Options) (*ZapLogger, error) {
	var sampling = zap.SamplingConfig{
		Initial:    100,
		Thereafter: 100,
	}

	var level = zap.InfoLevel
	if opts.Development == true {
		level = zap.DebugLevel
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	encoding := OutputModeConsole
	if opts.JsonOutputMode {
		encoding = OutputModeJson
	}

	// Build the zap logger
	zapLogger, err := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      opts.Development,
		Sampling:         &sampling,
		Encoding:         encoding,
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()
	if err != nil {
		return nil, err
	}
	coreLogger := zapLogger.WithOptions(zap.AddCallerSkip(opts.CallerSkip))
	return &ZapLogger{
		opts:  opts,
		core:  coreLogger,
		sugar: coreLogger.Sugar(),
	}, nil
}

func (l *ZapLogger) Clone(addedCallerSkip int) Logger {
	cp := *l
	cp.core = l.core.WithOptions(zap.AddCallerSkip(addedCallerSkip))
	cp.sugar = cp.core.Sugar()
	return &cp
}

func (l *ZapLogger) WithCtx(ctx context.Context) Logger {
	if l.opts.ContextExtractor == nil {
		return l
	}
	fields := l.opts.ContextExtractor(ctx)
	return l.WithField(fields...)
}

func (l *ZapLogger) WithField(fields ...field.Field) Logger {
	cp := *l
	cp.core = cp.core.With(fields...)
	cp.sugar = cp.core.Sugar()
	return &cp
}

func (l *ZapLogger) Info(args ...interface{}) {
	l.sugar.Info(args...)
}

func (l *ZapLogger) Infof(template string, args ...interface{}) {
	l.sugar.Infof(template, args...)
}

func (l *ZapLogger) Debug(args ...interface{}) {
	l.sugar.Debug(args...)
}

func (l *ZapLogger) Debugf(template string, args ...interface{}) {
	l.sugar.Debugf(template, args...)
}

func (l *ZapLogger) Warn(args ...interface{}) {
	l.sugar.Warn(args...)
}

func (l *ZapLogger) Warnf(template string, args ...interface{}) {
	l.sugar.Warnf(template, args...)
}

func (l *ZapLogger) Error(args ...interface{}) {
	l.sugar.Error(args...)
}

func (l *ZapLogger) Errorf(template string, args ...interface{}) {
	l.sugar.Errorf(template, args...)
}

func (l *ZapLogger) Fatal(args ...interface{}) {
	l.sugar.Fatal(args...)
}

func (l *ZapLogger) Fatalf(template string, args ...interface{}) {
	l.sugar.Fatalf(template, args...)
}

func (l *ZapLogger) Panic(args ...interface{}) {
	l.sugar.Panic(args...)
}

func (l *ZapLogger) Panicf(template string, args ...interface{}) {
	l.sugar.Panicf(template, args...)
}
