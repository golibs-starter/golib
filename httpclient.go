package golib

import (
	"fmt"
	"gitlab.id.vin/vincart/golib/config"
	"gitlab.id.vin/vincart/golib/web/client"
	"go.uber.org/fx"
)

func HttpClientOpt() fx.Option {
	return fx.Options(
		ProvideProps(client.NewHttpClientProperties),
		fx.Provide(NewContextualHttpClient),
	)
}

type ContextualHttpClientWrapper func(client.ContextualHttpClient) (client.ContextualHttpClient, error)

type HttpClientAutoConfigIn struct {
	fx.In
	AppProperties        *config.AppProperties
	HttpClientProperties *client.HttpClientProperties
	Wrappers             []ContextualHttpClientWrapper `group:"contextual_http_client_wrapper"`
}

// NewContextualHttpClient Initiate a client.ContextualHttpClient with
// configs are loaded automatically.
// Alternatively you can wrap the default client.ContextualHttpClient with
// one or more other client.ContextualHttpClient to customize the behavior.
// To do that, your provider have to return ContextualHttpClientWrapper.
// For example https://gitlab.id.vin/vincart/golib-security/-/blob/develop/httpclient.go
func NewContextualHttpClient(in HttpClientAutoConfigIn) (client.ContextualHttpClient, error) {
	// Create default http client
	defaultHttpClient, err := client.NewDefaultHttpClient(in.HttpClientProperties)
	if err != nil {
		return nil, fmt.Errorf("error when init default http client: [%v]", err)
	}

	// Wrap around by TraceableHttpClient by default
	var httpClient = client.NewTraceableHttpClient(defaultHttpClient, in.AppProperties)

	// Wrap around by other wrappers
	for _, wrapper := range in.Wrappers {
		httpClient, err = wrapper(httpClient)
		if err != nil {
			return nil, err
		}
	}

	return httpClient, nil
}
