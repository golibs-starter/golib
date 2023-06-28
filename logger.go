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

type NewLoggerOut struct {
	fx.Out
	Core log.Logger
	Web  log.Logger `name:"web_logger"`
}

func NewZapLogger(props *log.Properties) (NewLoggerOut, error) {
	out := NewLoggerOut{}
	// Create new logger instance
	logger, err := log.NewZapLogger(&log.Options{
		Development:      props.Development,
		JsonOutputMode:   props.JsonOutputMode,
		CallerSkip:       props.CallerSkip,
		ContextExtractor: webLog.ContextExtractor,
	})
	if err != nil {
		return out, errors.WithMessage(err, "init logger failed")
	}
	out.Core = logger
	log.ReplaceGlobal(logger)
	out.Web = logger.Clone(1)
	webLog.ReplaceGlobal(out.Web)
	return out, nil
}

type RegisterLoggerIn struct {
	fx.In
	Core log.Logger
	Web  log.Logger `name:"web_logger"`
}

func RegisterLogger(in RegisterLoggerIn) {
	// This is dummy invoker to make sure logger are produced by fx
}
