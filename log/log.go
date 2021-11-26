package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct {
	options     *Options
	coreLogger  *zap.Logger
	sugarLogger *zap.SugaredLogger
}

func NewLogger(options *Options) (*logger, error) {
	var sampling = zap.SamplingConfig{
		Initial:    100,
		Thereafter: 100,
	}

	var level = zap.InfoLevel
	if options.Development == true {
		level = zap.DebugLevel
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	encoding := OutputModeConsole
	if options.JsonOutputMode {
		encoding = OutputModeJson
	}

	// Build the zap logger
	zapLogger, err := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      options.Development,
		Sampling:         &sampling,
		Encoding:         encoding,
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()
	if err != nil {
		return nil, err
	}
	coreLogger := zapLogger.WithOptions(zap.AddCallerSkip(options.CallerSkip))
	return &logger{
		options:     options,
		coreLogger:  coreLogger,
		sugarLogger: coreLogger.Sugar(),
	}, nil
}

func (l *logger) Clone(options ...OptionFunc) Logger {
	cp := *l
	newOpt := *l.options
	cp.options = &newOpt
	for _, opt := range options {
		opt(cp.options)
	}
	cp.coreLogger = l.coreLogger.WithOptions(zap.AddCallerSkip(cp.options.CallerSkip))
	cp.sugarLogger = cp.coreLogger.Sugar()
	return &cp
}

func (l *logger) logw(level zapcore.Level, keysAndValues []interface{}, msgFormat string, args ...interface{}) {
	msg := msgFormat
	if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(msgFormat, args...)
	}
	switch level {
	case zapcore.DebugLevel:
		l.sugarLogger.Debugw(msg, keysAndValues...)
		break
	case zapcore.InfoLevel:
		l.sugarLogger.Infow(msg, keysAndValues...)
		break
	case zapcore.WarnLevel:
		l.sugarLogger.Warnw(msg, keysAndValues...)
		break
	case zapcore.ErrorLevel:
		l.sugarLogger.Errorw(msg, keysAndValues...)
		break
	case zapcore.FatalLevel:
		l.sugarLogger.Fatalw(msg, keysAndValues...)
		break
	}
}

func (l *logger) Info(args ...interface{}) {
	l.sugarLogger.Info(args...)
}

func (l *logger) Infof(msgFormat string, args ...interface{}) {
	l.sugarLogger.Infof(msgFormat, args...)
}

func (l *logger) Infow(keysAndValues []interface{}, msgFormat string, args ...interface{}) {
	l.logw(zapcore.InfoLevel, keysAndValues, msgFormat, args...)
}

func (l *logger) Debug(args ...interface{}) {
	l.sugarLogger.Debug(args...)
}

func (l *logger) Debugf(msgFormat string, args ...interface{}) {
	l.sugarLogger.Debugf(msgFormat, args...)
}

func (l *logger) Debugw(keysAndValues []interface{}, msgFormat string, args ...interface{}) {
	l.logw(zapcore.DebugLevel, keysAndValues, msgFormat, args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.sugarLogger.Warn(args...)
}

func (l *logger) Warnf(msgFormat string, args ...interface{}) {
	l.sugarLogger.Warnf(msgFormat, args...)
}

func (l *logger) Warnw(keysAndValues []interface{}, msgFormat string, args ...interface{}) {
	l.logw(zapcore.WarnLevel, keysAndValues, msgFormat, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.sugarLogger.Error(args...)
}

func (l *logger) Errorf(msgFormat string, args ...interface{}) {
	l.sugarLogger.Errorf(msgFormat, args...)
}

func (l *logger) Errorw(keysAndValues []interface{}, msgFormat string, args ...interface{}) {
	l.logw(zapcore.ErrorLevel, keysAndValues, msgFormat, args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.sugarLogger.Fatal(args...)
}

func (l *logger) Fatalf(msgFormat string, args ...interface{}) {
	l.sugarLogger.Fatalf(msgFormat, args...)
}

func (l *logger) Fatalw(keysAndValues []interface{}, msgFormat string, args ...interface{}) {
	l.logw(zapcore.FatalLevel, keysAndValues, msgFormat, args...)
}
