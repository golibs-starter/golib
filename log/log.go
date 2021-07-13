package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Options struct {
	// Development puts the logger in development mode, which changes the
	// behavior of DPanicLevel and takes stack traces more liberally.
	Development bool

	// Enable json output mode
	JsonOutputMode bool

	// Skip number of callers before show caller
	CallerSkip int
}

const (
	OutputModeJson    = "json"
	OutputModeConsole = "console"
)

type logger struct {
	options *Options
}

func NewLogger(options *Options) (*logger, error) {
	var sampling = zap.SamplingConfig{
		Initial:    100,
		Thereafter: 100,
	}

	// Default behavior for the logger
	var level zapcore.Level
	var encoderConfig zapcore.EncoderConfig
	if options.Development == true {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		level = zap.DebugLevel
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
		level = zap.InfoLevel
	}

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

	zapOptions := append(make([]zap.Option, 0), zap.AddCallerSkip(options.CallerSkip))
	zap.ReplaceGlobals(zapLogger.WithOptions(zapOptions...))
	return &logger{options: options}, nil
}

func (l *logger) logw(level zapcore.Level, keysAndValues []interface{}, msgFormat string, args ...interface{}) {
	msg := msgFormat
	if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(msgFormat, args...)
	}
	switch level {
	case zapcore.DebugLevel:
		zap.S().Debugw(msg, keysAndValues...)
		break
	case zapcore.InfoLevel:
		zap.S().Infow(msg, keysAndValues...)
		break
	case zapcore.WarnLevel:
		zap.S().Warnw(msg, keysAndValues...)
		break
	case zapcore.ErrorLevel:
		zap.S().Errorw(msg, keysAndValues...)
		break
	case zapcore.FatalLevel:
		zap.S().Fatalw(msg, keysAndValues...)
		break
	}
}

func (l *logger) Info(args ...interface{}) {
	zap.S().Info(args...)
}

func (l *logger) Infof(msgFormat string, args ...interface{}) {
	zap.S().Infof(msgFormat, args...)
}

func (l *logger) Infow(keysAndValues []interface{}, msgFormat string, args ...interface{}) {
	l.logw(zapcore.InfoLevel, keysAndValues, msgFormat, args...)
}

func (l *logger) Debug(args ...interface{}) {
	zap.S().Debug(args...)
}

func (l *logger) Debugf(msgFormat string, args ...interface{}) {
	zap.S().Debugf(msgFormat, args...)
}

func (l *logger) Debugw(keysAndValues []interface{}, msgFormat string, args ...interface{}) {
	l.logw(zapcore.DebugLevel, keysAndValues, msgFormat, args...)
}

func (l *logger) Warn(args ...interface{}) {
	zap.S().Warn(args...)
}

func (l *logger) Warnf(msgFormat string, args ...interface{}) {
	zap.S().Warnf(msgFormat, args...)
}

func (l *logger) Warnw(keysAndValues []interface{}, msgFormat string, args ...interface{}) {
	l.logw(zapcore.WarnLevel, keysAndValues, msgFormat, args...)
}

func (l *logger) Error(args ...interface{}) {
	zap.S().Error(args...)
}

func (l *logger) Errorf(msgFormat string, args ...interface{}) {
	zap.S().Errorf(msgFormat, args...)
}

func (l *logger) Errorw(keysAndValues []interface{}, msgFormat string, args ...interface{}) {
	l.logw(zapcore.ErrorLevel, keysAndValues, msgFormat, args...)
}

func (l *logger) Fatal(args ...interface{}) {
	zap.S().Fatal(args...)
}

func (l *logger) Fatalf(msgFormat string, args ...interface{}) {
	zap.S().Fatalf(msgFormat, args...)
}

func (l *logger) Fatalw(keysAndValues []interface{}, msgFormat string, args ...interface{}) {
	l.logw(zapcore.FatalLevel, keysAndValues, msgFormat, args...)
}
