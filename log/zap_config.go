package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func buildZapLoggerConfig(opts *Options) zap.Config {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        opts.FieldKeyMap[FieldKeyTime],
		LevelKey:       opts.FieldKeyMap[FieldKeyLevel],
		NameKey:        opts.FieldKeyMap[FieldKeyName],
		CallerKey:      opts.FieldKeyMap[FieldKeyCaller],
		FunctionKey:    opts.FieldKeyMap[FieldKeyFunc],
		MessageKey:     opts.FieldKeyMap[FieldKeyMsg],
		StacktraceKey:  opts.FieldKeyMap[FieldKeyStacktrace],
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var level = zap.InfoLevel
	if opts.Development {
		level = zap.DebugLevel
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	if opts.LogLevel != "" {
		lv := zap.NewAtomicLevel()
		if err := lv.UnmarshalText([]byte(opts.LogLevel)); err == nil {
			level = lv.Level()
		}
	}

	encoding := OutputModeConsole
	if opts.JsonOutputMode {
		encoding = OutputModeJson
	}

	// Build the zap logger
	return zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       opts.Development,
		DisableCaller:     opts.DisableCaller,
		DisableStacktrace: opts.DisableStacktrace,
		Encoding:          encoding,
		EncoderConfig:     encoderConfig,
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
	}
}
