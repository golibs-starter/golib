package golib

import (
	"gitlab.id.vin/vincart/golib/config"
	"gitlab.id.vin/vincart/golib/web/client"
	"gitlab.id.vin/vincart/golib/web/log"
)

type Properties struct {
	Application *config.ApplicationProperties
	Logging     *log.LoggingProperties
	HttpClient  *client.HttpClientProperties
}

type ConfigOption func(option *config.Option)

func SetActiveProfiles(activeProfiles []string) ConfigOption {
	return func(option *config.Option) {
		option.ActiveProfiles = activeProfiles
	}
}

func SetConfigPaths(configPaths []string) ConfigOption {
	return func(option *config.Option) {
		option.ConfigPaths = configPaths
	}
}

// SetConfigFormat accept yaml, json values
func SetConfigFormat(configFormat string) ConfigOption {
	return func(option *config.Option) {
		option.ConfigFormat = configFormat
	}
}

func WithConfigProperties(options ...ConfigOption) Module {
	return func(app *App) {
		option := new(config.Option)
		for _, optFunc := range options {
			optFunc(option)
		}
		app.ConfigLoader = config.NewLoader(*option, nil)
		app.Properties = &Properties{}

		// Bind application properties
		app.Properties.Application = &config.ApplicationProperties{}
		app.ConfigLoader.Bind(app.Properties.Application)
	}
}
