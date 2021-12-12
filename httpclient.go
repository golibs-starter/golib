package golib

import (
	"fmt"
	"gitlab.com/golibs-starter/golib/config"
	"gitlab.com/golibs-starter/golib/web/client"
	"go.uber.org/fx"
	"net/http"
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

type HttpClientAutoConfigOut struct {
	fx.Out
	Base                 *http.Client
	HttpClient           client.HttpClient
	ContextualHttpClient client.ContextualHttpClient
}

// NewContextualHttpClient Initiate a client.ContextualHttpClient with
// configs are loaded automatically.
// Alternatively you can wrap the default client.ContextualHttpClient with
// one or more other client.ContextualHttpClient to customize the behavior.
// To do that, your provider have to return ContextualHttpClientWrapper.
// For example https://gitlab.com/golibs-starter/golib-security/-/blob/develop/httpclient.go
func NewContextualHttpClient(in HttpClientAutoConfigIn) (HttpClientAutoConfigOut, error) {
	out := HttpClientAutoConfigOut{Base: &http.Client{}}

	// Create default http client
	defaultHttpClient, err := client.NewDefaultHttpClient(out.Base, in.HttpClientProperties)
	if err != nil {
		return out, fmt.Errorf("error when init default http client: [%v]", err)
	}
	out.HttpClient = defaultHttpClient

	// Wrap around by TraceableHttpClient by default
	var httpClient = client.NewTraceableHttpClient(defaultHttpClient, in.AppProperties)

	// Wrap around by other wrappers
	for _, wrapper := range in.Wrappers {
		httpClient, err = wrapper(httpClient)
		if err != nil {
			return out, err
		}
	}
	out.ContextualHttpClient = httpClient
	return out, nil
}
