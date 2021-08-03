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
	Properties   *Properties
	ConfigLoader config.Loader
	Logger       log.Logger
	Publisher    pubsub.Publisher
	HttpClient   client.ContextualHttpClient
	middleware   []func(next http.Handler) http.Handler
}

func New(modules ...Module) *App {
	app := new(App)
	app.Properties = new(Properties)
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

func (a App) Port() int {
	return a.Properties.Application.Port
}

func (a App) Name() string {
	return a.Properties.Application.Name
}
