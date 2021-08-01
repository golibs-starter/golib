package golib

import (
	"gitlab.id.vin/vincart/golib/config"
	"gitlab.id.vin/vincart/golib/utils"
	"gitlab.id.vin/vincart/golib/web/client"
	"gitlab.id.vin/vincart/golib/web/log"
	"os"
)

type Properties struct {
	Application *config.ApplicationProperties
	Logging     *log.LoggingProperties
	HttpClient  *client.HttpClientProperties
}

type ConfigOption func(option *config.Option)

func OptActiveProfiles(activeProfiles []string) ConfigOption {
	return func(option *config.Option) {
		option.ActiveProfiles = activeProfiles
	}
}

func OptConfigPaths(configPaths []string) ConfigOption {
	return func(option *config.Option) {
		option.ConfigPaths = configPaths
	}
}

// OptConfigFormat accept yaml, json values
func OptConfigFormat(configFormat string) ConfigOption {
	return func(option *config.Option) {
		option.ConfigFormat = configFormat
	}
}

func OptConfigFromEnv() ConfigOption {
	return func(option *config.Option) {
		option.ActiveProfiles = utils.SliceFromCommaString(os.Getenv("ENV"))
		option.ConfigPaths = utils.SliceFromCommaString(os.Getenv("CONFIG_PATHS"))
		option.ConfigFormat = os.Getenv("CONFIG_FORMAT")
	}
}

func WithProperties(options ...ConfigOption) Module {
	if len(options) == 0 {
		options = append(options, OptConfigFromEnv())
	}
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
