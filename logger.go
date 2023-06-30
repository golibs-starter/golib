package golib

import (
	"github.com/pkg/errors"
	"gitlab.com/golibs-starter/golib/log"
	webLog "gitlab.com/golibs-starter/golib/web/log"
	"go.uber.org/fx"
)

func LoggingOpt() fx.Option {
	return fx.Options(
		ProvideProps(log.NewProperties),
		fx.Provide(NewZapLogger),
		fx.Invoke(RegisterLogger),
	)
}

func NewZapLogger(props *log.Properties) (log.Logger, error) {
	// Create new logger instance
	logger, err := log.NewZapLogger(&log.Options{
		Development:      props.Development,
		JsonOutputMode:   props.JsonOutputMode,
		CallerSkip:       props.CallerSkip,
		ContextExtractor: webLog.ContextExtractor,
	})
	if err != nil {
		return nil, errors.WithMessage(err, "init logger failed")
	}
	log.ReplaceGlobal(logger)
	webLog.ReplaceGlobal(logger.Clone(1))
	return logger, nil
}

type RegisterLoggerIn struct {
	fx.In
	Core log.Logger
}

func RegisterLogger(in RegisterLoggerIn) {
	// This is dummy invoker to make sure logger are produced by fx
}
