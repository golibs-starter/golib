package golib

import (
	"fmt"
	"gitlab.id.vin/vincart/golib/config"
	"gitlab.id.vin/vincart/golib/web/client"
)

type ContextualHttpClientWrapper func(client.ContextualHttpClient) client.ContextualHttpClient
type ContextualHttpClientWrappers []ContextualHttpClientWrapper

func NewHttpClientAutoConfig(
	loader config.Loader,
	appProps *config.AppProperties,
	wrappers ContextualHttpClientWrappers,
) (client.ContextualHttpClient, *client.HttpClientProperties, error) {
	props := client.NewHttpClientProperties(loader)
	httpClient, err := NewContextualHttpClient(appProps, props, wrappers)
	if err != nil {
		return nil, nil, err
	}
	return httpClient, props, nil
}

func NewContextualHttpClient(
	appProps *config.AppProperties,
	httpClientProps *client.HttpClientProperties,
	wrappers ContextualHttpClientWrappers,
) (client.ContextualHttpClient, error) {
	// Create default http client
	defaultHttpClient, err := client.NewDefaultHttpClient(httpClientProps)
	if err != nil {
		return nil, fmt.Errorf("error when init default http client: [%v]", err)
	}

	// Wrap around by TraceableHttpClient by default
	var httpClient = client.NewTraceableHttpClient(
		defaultHttpClient, appProps,
	)

	// Wrap around by other wrappers
	for _, wrapper := range wrappers {
		httpClient = wrapper(httpClient)
	}
	return httpClient, nil
}
