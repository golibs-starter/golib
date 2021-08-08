package golib

import (
	"fmt"
	"gitlab.id.vin/vincart/golib/config"
	"gitlab.id.vin/vincart/golib/web/client"
	"go.uber.org/fx"
)

type ContextualHttpClientWrapper func(client.ContextualHttpClient) (client.ContextualHttpClient, error)

type HttpClientAutoConfigIn struct {
	fx.In
	ConfigLoader  config.Loader
	AppProperties *config.AppProperties
	Wrappers      []ContextualHttpClientWrapper `group:"contextual_http_client_wrapper"`
}

type HttpClientAutoConfigOut struct {
	fx.Out
	Properties *client.HttpClientProperties
	HttpClient client.ContextualHttpClient
}

// NewHttpClientAutoConfig Initiate a client.ContextualHttpClient with
// configs are loaded automatically.
// Alternatively you can wrap the default client.ContextualHttpClient with
// one or more other client.ContextualHttpClient to customize the behavior.
// To do that, your provider have to return ContextualHttpClientWrapper.
// For example https://gitlab.id.vin/vincart/golib-security/-/blob/develop/httpclient.go
func NewHttpClientAutoConfig(in HttpClientAutoConfigIn) (HttpClientAutoConfigOut, error) {
	out := HttpClientAutoConfigOut{}
	properties, err := client.NewHttpClientProperties(in.ConfigLoader)
	if err != nil {
		return out, err
	}

	// Create default http client
	defaultHttpClient, err := client.NewDefaultHttpClient(properties)
	if err != nil {
		return out, fmt.Errorf("error when init default http client: [%v]", err)
	}

	// Wrap around by TraceableHttpClient by default
	var httpClient = client.NewTraceableHttpClient(defaultHttpClient, in.AppProperties)

	// Wrap around by other wrappers
	for _, wrapper := range in.Wrappers {
		httpClient, err = wrapper(httpClient)
		if err != nil {
			return out, err
		}
	}

	out.Properties = properties
	out.HttpClient = httpClient
	return out, nil
}
