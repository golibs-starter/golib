package golib

import (
	"gitlab.id.vin/vincart/golib/config"
	"gitlab.id.vin/vincart/golib/pubsub"
	"gitlab.id.vin/vincart/golib/web/client"
	"gitlab.id.vin/vincart/golib/web/middleware"
	"net/http"
)

type Options struct {
	EventMapping map[pubsub.Event][]pubsub.Subscriber
}

type App struct {
	middleware []func(next http.Handler) http.Handler
}

func (a App) Middleware() []func(next http.Handler) http.Handler {
	return a.middleware
}

func Init(options Options) *App {
	httpClientProperties := new(client.HttpClientConfig)
	InitConfig(config.Option{}, []config.Properties{httpClientProperties})

	logger := InitLogger()
	InitEventBus(options.EventMapping, logger)
	return &App{
		middleware: defaultMiddleware(),
	}
}

func defaultMiddleware() []func(next http.Handler) http.Handler {
	return []func(next http.Handler) http.Handler{
		middleware.AdvancedResponseWriter(),
		middleware.RequestContext(),
		middleware.CorrelationId(),
	}
}
