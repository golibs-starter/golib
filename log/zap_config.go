package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func buildZapLoggerConfig(opts *Options) zap.Config {
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
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: opts.Development,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         encoding,
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}
