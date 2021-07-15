package golib

import (
	"fmt"
	"gitlab.id.vin/vincart/golib/web/client"
)

func WithHttpClientAutoConfig() Module {
	return func(app *App) {
		app.Properties.HttpClientProperties = client.HttpClientProperties{}
		app.Loader.Bind(&app.Properties.HttpClientProperties)
		httpClient, err := client.NewDefaultHttpClient(&app.Properties.HttpClientProperties)
		if err != nil {
			panic(fmt.Sprintf("Error when init default http client: [%v]", err))
		}
		app.HttpClient = httpClient
	}
}
