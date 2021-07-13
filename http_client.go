package golib

import "gitlab.id.vin/vincart/golib/web/client"

func WithHttpClientAutoConfig() Module {
	return func(app *App) {
		app.Properties.HttpClientProperties = client.HttpClientProperties{}
		app.Loader.Bind(&app.Properties.HttpClientProperties)
	}
}
