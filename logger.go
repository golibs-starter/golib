package golib

import (
	"fmt"
	"gitlab.id.vin/vincart/golib/log"
	webLog "gitlab.id.vin/vincart/golib/web/log"
)

func WithLoggerAutoConfig() Module {
	return func(app *App) {
		// Bind logging properties
		app.Properties.Logging = &webLog.LoggingProperties{}
		app.ConfigLoader.Bind(app.Properties.Logging)

		// Create new logger instance
		logger, err := log.NewLogger(&log.Options{
			Development:    app.Properties.Logging.Development,
			JsonOutputMode: app.Properties.Logging.JsonOutputMode,
			CallerSkip:     app.Properties.Logging.CallerSkip,
		})
		if err != nil {
			panic(fmt.Sprintf("Error when init logger: [%v]", err))
		}
		log.RegisterGlobal(logger)
		app.Logger = logger
	}
}
