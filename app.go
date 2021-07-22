package golib

import (
	"gitlab.id.vin/vincart/golib/config"
	"gitlab.id.vin/vincart/golib/log"
	"gitlab.id.vin/vincart/golib/pubsub"
	"gitlab.id.vin/vincart/golib/web/client"
	"gitlab.id.vin/vincart/golib/web/middleware"
	"net/http"
)

type Module func(app *App)

type App struct {
	middleware   []func(next http.Handler) http.Handler
	Properties   *Properties
	ConfigLoader config.Loader
	Logger       log.Logger
	Publisher    pubsub.Publisher
	HttpClient   client.ContextualHttpClient
}

func New(modules ...Module) *App {
	app := new(App)
	app.AddMiddleware(
		middleware.AdvancedResponseWriter(),
		middleware.RequestContext(),
		middleware.CorrelationId(),
	)
	for _, module := range modules {
		module(app)
	}
	return app
}

func (a *App) AddMiddleware(middleware ...func(next http.Handler) http.Handler) {
	a.middleware = append(a.middleware, middleware...)
}

func (a App) Middleware() []func(next http.Handler) http.Handler {
	return a.middleware
}
