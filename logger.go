package golib

import (
	"fmt"
	"gitlab.id.vin/vincart/golib/config"
	"gitlab.id.vin/vincart/golib/log"
	webLog "gitlab.id.vin/vincart/golib/web/log"
)

func NewLoggerAutoConfig(loader config.Loader) (log.Logger, *webLog.LoggingProperties, error) {
	props := webLog.NewLoggingProperties(loader)
	logger, err := NewLogger(props)
	if err != nil {
		return nil, nil, err
	}
	return logger, props, nil
}

func RegisterLoggerAutoConfig(logger log.Logger) {
	log.ReplaceGlobal(logger)
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
