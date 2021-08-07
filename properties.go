package golib

import (
	"gitlab.id.vin/vincart/golib/config"
	"gitlab.id.vin/vincart/golib/log"
	"gitlab.id.vin/vincart/golib/utils"
	"os"
)

type Option func(option *config.Option)

func WithActiveProfiles(activeProfiles []string) Option {
	return func(option *config.Option) {
		option.ActiveProfiles = activeProfiles
	}
}

func WithPaths(configPaths []string) Option {
	return func(option *config.Option) {
		option.ConfigPaths = configPaths
	}
}

// WithFormat accept yaml, json values
func WithFormat(configFormat string) Option {
	return func(option *config.Option) {
		option.ConfigFormat = configFormat
	}
}

func WithEnvironmentOption() Option {
	return func(option *config.Option) {
		option.ActiveProfiles = utils.SliceFromCommaString(os.Getenv("APPLICATION_ENV"))
		option.ConfigPaths = utils.SliceFromCommaString(os.Getenv("APPLICATION_CONFIG_PATHS"))
		option.ConfigFormat = os.Getenv("APPLICATION_CONFIG_FORMAT")
	}
}

func NewPropertiesAutoLoad(options ...Option) (config.Loader, *config.AppProperties, error) {
	opts := []Option{WithEnvironmentOption()}
	opts = append(opts, options...)
	option := new(config.Option)
	for _, optFunc := range opts {
		optFunc(option)
	}
	configLoader, err := config.NewLoader(*option, log.Debugf)
	if err != nil {
		return nil, nil, err
	}
	props, err := config.NewApplicationProperties(configLoader)
	if err != nil {
		return nil, nil, err
	}
	return configLoader, props, nil
}
