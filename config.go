package golib

import (
	"gitlab.id.vin/vincart/golib/config"
	"gitlab.id.vin/vincart/golib/web/client"
)

type Properties struct {
	Application *config.ApplicationProperties
	HttpClient  *client.HttpClientProperties
}

func WithConfigLoader(option config.Option) Module {
	return func(app *App) {
		app.ConfigLoader = config.NewLoader(option, nil)
		app.Properties = &Properties{}

		// Bind application properties
		app.Properties.Application = &config.ApplicationProperties{}
		app.ConfigLoader.Bind(app.Properties.Application)
	}
}
