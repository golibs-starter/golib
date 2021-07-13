package golib

import (
	"gitlab.id.vin/vincart/golib/config"
	"gitlab.id.vin/vincart/golib/web/client"
)

type Properties struct {
	HttpClientProperties client.HttpClientProperties
}

func WithConfigLoader(option config.Option) Module {
	return func(app *App) {
		app.Loader = config.NewLoader(option, nil)
		app.Properties = &Properties{}
	}
}
