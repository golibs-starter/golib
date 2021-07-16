package golib

import (
	"fmt"
	"gitlab.id.vin/vincart/golib/web/client"
)

type HttpClientWrapper func(client.HttpClient) client.HttpClient

func WithHttpClientAutoConfig(wrappers ...HttpClientWrapper) Module {
	return func(app *App) {
		app.Properties.HttpClientProperties = client.HttpClientProperties{}
		app.Loader.Bind(&app.Properties.HttpClientProperties)
		defaultHttpClient, err := client.NewDefaultHttpClient(&app.Properties.HttpClientProperties)
		if err != nil {
			panic(fmt.Sprintf("Error when init default http client: [%v]", err))
		}
		var httpClient client.HttpClient = defaultHttpClient
		for _, wrapper := range wrappers {
			httpClient = wrapper(httpClient)
		}
		app.HttpClient = httpClient
	}
}
