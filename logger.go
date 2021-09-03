package golib

import (
	"fmt"
	"gitlab.id.vin/vincart/golib/log"
	webLog "gitlab.id.vin/vincart/golib/web/log"
	"go.uber.org/fx"
)

func LoggingOpt() fx.Option {
	return fx.Options(
		EnablePropsAutoload(new(webLog.LoggingProperties)),
		fx.Provide(webLog.NewLoggingProperties),
		fx.Provide(NewLogger),
		fx.Invoke(RegisterLogger),
	)
}

func NewLogger(props *webLog.LoggingProperties) (log.Logger, error) {
	// Create new logger instance
	logger, err := log.NewLogger(&log.Options{
		Development:    props.Development,
		JsonOutputMode: props.JsonOutputMode,
		CallerSkip:     props.CallerSkip,
	})
	if err != nil {
		return nil, fmt.Errorf("error when init logger: [%v]", err)
	}
	return logger, nil
}

func RegisterLogger(logger log.Logger) {
	log.ReplaceGlobal(logger)
}
