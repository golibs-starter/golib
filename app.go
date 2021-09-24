package golib

import (
	"gitlab.id.vin/vincart/golib/config"
	"gitlab.id.vin/vincart/golib/web/middleware"
	"go.uber.org/fx"
	"net/http"
)

func AppOpt() fx.Option {
	return fx.Options(
		ProvideProps(config.NewAppProperties),
		fx.Provide(New),
	)
}

func New(props *config.AppProperties) *App {
	app := App{Properties: props}
	app.AddHandler(
		middleware.AdvancedResponseWriter(),
		middleware.RequestContext(),
		middleware.CorrelationId(),
	)
	return &app
}

type App struct {
	Properties *config.AppProperties
	handlers   []func(next http.Handler) http.Handler
}

func (a App) Name() string {
	return a.Properties.Name
}

func (a App) Port() int {
	return a.Properties.Port
}

func (a App) Path() string {
	return a.Properties.Path
}

func (a App) Handlers() []func(next http.Handler) http.Handler {
	return a.handlers
}

func (a *App) AddHandler(handlers ...func(next http.Handler) http.Handler) {
	a.handlers = append(a.handlers, handlers...)
}
