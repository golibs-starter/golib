package golib

import (
	"github.com/golibs-starter/golib/log"
	webLog "github.com/golibs-starter/golib/web/log"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

func LoggingOpt() fx.Option {
	return fx.Options(
		ProvideProps(log.NewProperties),
		RegisterLogContextExtractor(webLog.ContextExtractor),
		fx.Provide(NewZapLogger),
		fx.Invoke(RegisterLogger),
	)
}

func RegisterLogContextExtractor(extractor log.ContextExtractor) fx.Option {
	return fx.Provide(fx.Annotated{
		Group:  "log_context_extractor",
		Target: func() log.ContextExtractor { return extractor },
	})
}

type ZapLoggerIn struct {
	fx.In
	Props             *log.Properties
	ContextExtractors []log.ContextExtractor `group:"log_context_extractor"`
}

func NewZapLogger(in ZapLoggerIn) (log.Logger, error) {
	// Create new logger instance
	logger, err := log.NewZapLogger(&log.Options{
		Development:       in.Props.Development,
		LogLevel:          in.Props.LogLevel,
		JsonOutputMode:    in.Props.JsonOutputMode,
		DisableCaller:     in.Props.DisableCaller,
		DisableStacktrace: in.Props.DisableStacktrace,
		CallerSkip:        in.Props.CallerSkip,
		ContextExtractors: in.ContextExtractors,
	})
	if err != nil {
		return nil, errors.WithMessage(err, "init logger failed")
	}
	log.ReplaceGlobal(logger)
	webLog.ReplaceGlobal(logger.Clone(1))
	return logger, nil
}

func RegisterLogger(logger log.Logger) {
	// This is dummy invoker to make sure logger are produced by fx
}
