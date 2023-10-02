package golib

import (
	"github.com/golibs-starter/golib/config"
	"github.com/golibs-starter/golib/web/client"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

func HttpClientOpt() fx.Option {
	return fx.Options(
		ProvideProps(client.NewHttpClientProperties),
		fx.Provide(client.NewNativeHttpClient),
		fx.Provide(client.NewDefaultHttpClient),
		fx.Provide(NewContextualHttpClient),
	)
}

type ContextualHttpClientWrapper func(client.ContextualHttpClient) (client.ContextualHttpClient, error)

type ContextualHttpClientIn struct {
	fx.In
	AppProperties *config.AppProperties
	HttpClient    client.HttpClient
	Wrappers      []ContextualHttpClientWrapper `group:"contextual_http_client_wrapper"`
}

// NewContextualHttpClient Initiate a client.ContextualHttpClient with
// configs are loaded automatically.
// Alternatively you can wrap the default client.ContextualHttpClient with
// one or more other client.ContextualHttpClient to customize the behavior.
// To do that, your provider have to return ContextualHttpClientWrapper.
// For example https://github.com/golibs-starter/golib-security/-/blob/develop/httpclient.go
func NewContextualHttpClient(in ContextualHttpClientIn) (client.ContextualHttpClient, error) {
	// Wrap around by TraceableHttpClient by default
	var httpClient = client.NewTraceableHttpClient(in.HttpClient, in.AppProperties)

	// Wrap around by other wrappers
	var err error
	for _, wrapper := range in.Wrappers {
		httpClient, err = wrapper(httpClient)
		if err != nil {
			return nil, errors.WithMessage(err, "wrap contextual http client failed")
		}
	}

	return httpClient, nil
}
