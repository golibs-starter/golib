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
	core := zapLogger.WithOptions(zap.AddCallerSkip(opts.CallerSkip))
	return &ZapLogger{
		opts:  opts,
		core:  core,
		sugar: core.Sugar(),
	}, nil
}

func (l *ZapLogger) Clone(addedCallerSkip int, fields ...field.Field) *ZapLogger {
	cp := *l
	cp.core = l.core.WithOptions(zap.AddCallerSkip(addedCallerSkip), zap.Fields(fields...))
	cp.sugar = cp.core.Sugar()
	return &cp
}

func (l *ZapLogger) WithCtx(ctx context.Context, additionalFields ...field.Field) Logger {
	fields := additionalFields
	if l.opts.ContextExtractor != nil {
		fields = append(fields, l.opts.ContextExtractor(ctx)...)
	}
	return l.Clone(0, fields...)
}

func (l *ZapLogger) WithField(fields ...field.Field) Logger {
	cp := l.Clone(0, fields...)
	return cp
}

func (l *ZapLogger) prepareArgs(args ...interface{}) (StdLogger, []interface{}) {
	if len(args) > 0 {
		if ctx, ok := args[0].(context.Context); ok {
			return l.Clone(1).WithCtx(ctx), args[1:]
		}
	}
	return l.sugar, args
}

func (l *ZapLogger) Info(args ...interface{}) {
	stdLogger, args := l.prepareArgs(args...)
	stdLogger.Info(args...)
}

func (l *ZapLogger) Infof(template string, args ...interface{}) {
	l.sugar.Infof(template, args...)
}

func (l *ZapLogger) Debug(args ...interface{}) {
	stdLogger, args := l.prepareArgs(args...)
	stdLogger.Debug(args...)
}

func (l *ZapLogger) Debugf(template string, args ...interface{}) {
	l.sugar.Debugf(template, args...)
}

func (l *ZapLogger) Warn(args ...interface{}) {
	stdLogger, args := l.prepareArgs(args...)
	stdLogger.Warn(args...)
}

func (l *ZapLogger) Warnf(template string, args ...interface{}) {
	l.sugar.Warnf(template, args...)
}

func (l *ZapLogger) Error(args ...interface{}) {
	stdLogger, args := l.prepareArgs(args...)
	stdLogger.Error(args...)
}

func (l *ZapLogger) Errorf(template string, args ...interface{}) {
	l.sugar.Errorf(template, args...)
}

func (l *ZapLogger) Fatal(args ...interface{}) {
	stdLogger, args := l.prepareArgs(args...)
	stdLogger.Fatal(args...)
}

func (l *ZapLogger) Fatalf(template string, args ...interface{}) {
	l.sugar.Fatalf(template, args...)
}

func (l *ZapLogger) Panic(args ...interface{}) {
	stdLogger, args := l.prepareArgs(args...)
	stdLogger.Panic(args...)
}

func (l *ZapLogger) Panicf(template string, args ...interface{}) {
	l.sugar.Panicf(template, args...)
}
