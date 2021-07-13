package golib

import (
	"fmt"
	"gitlab.id.vin/vincart/golib/log"
)

func WithLogger() Module {
	return func(app *App) {
		logger, err := log.NewLogger(&log.Options{
			Development:    true,
			JsonOutputMode: true,
			CallerSkip:     2,
		})
		if err != nil {
			panic(fmt.Sprintf("Error when init logger: [%v]", err))
		}
		log.RegisterGlobal(logger)
		app.Logger = logger
	}
}
