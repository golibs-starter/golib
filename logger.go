package golib

import (
	"fmt"
	"gitlab.com/golibs-starter/golib/log"
	webLog "gitlab.com/golibs-starter/golib/web/log"
	"go.uber.org/fx"
)

func LoggingOpt() fx.Option {
	return fx.Options(
		ProvideProps(log.NewProperties),
		fx.Provide(NewLogger),
		fx.Invoke(RegisterLogger),
	)
}

type NewLoggerOut struct {
	fx.Out
	Core log.Logger
	Web  log.Logger `name:"web_logger"`
}

func NewLogger(props *log.Properties) (NewLoggerOut, error) {
	out := NewLoggerOut{}
	// Create new logger instance
	logger, err := log.NewLogger(&log.Options{
		Development:    props.Development,
		JsonOutputMode: props.JsonOutputMode,
		CallerSkip:     props.CallerSkip,
	})
	if err != nil {
		return out, fmt.Errorf("error when init logger: [%v]", err)
	}
	out.Core = logger
	out.Web = logger.Clone(log.AddCallerSkip(1))
	return out, nil
}

type RegisterLoggerIn struct {
	fx.In
	Core log.Logger
	Web  log.Logger `name:"web_logger"`
}

func RegisterLogger(in RegisterLoggerIn) {
	log.ReplaceGlobal(in.Core)
	webLog.ReplaceGlobal(in.Web)
}
