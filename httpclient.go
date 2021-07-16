package golib

import (
	"fmt"
	"gitlab.id.vin/vincart/golib/web/client"
)

type ContextualHttpClientWrapper func(httpClient client.ContextualHttpClient) client.ContextualHttpClient

func WithHttpClientAutoConfig(wrappers ...ContextualHttpClientWrapper) Module {
	return func(app *App) {
		// Bind http client properties
		app.Properties.HttpClient = &client.HttpClientProperties{}
		app.ConfigLoader.Bind(app.Properties.HttpClient)

		// Create default http client
		defaultHttpClient, err := client.NewDefaultHttpClient(app.Properties.HttpClient)
		if err != nil {
			panic(fmt.Sprintf("Error when init default http client: [%v]", err))
		}
		var httpClient client.ContextualHttpClient = client.NewTraceableHttpClient(
			defaultHttpClient, app.Properties.Application,
		)
		for _, wrapper := range wrappers {
			httpClient = wrapper(httpClient)
		}
		app.HttpClient = httpClient
	}
}
