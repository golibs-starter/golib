package golib

import (
	"gitlab.id.vin/vincart/golib/config"
	"gitlab.id.vin/vincart/golib/log"
	"gitlab.id.vin/vincart/golib/utils"
	"go.uber.org/fx"
	"os"
)

func PropertiesAutoConfig() fx.Option {
	return fx.Provide(NewPropertiesLoader)
}

type EnablePropsAutoloadOut struct {
	fx.Out
	Properties config.Properties `group:"properties"`
}

func EnablePropsAutoload(props config.Properties) fx.Option {
	return fx.Provide(func() EnablePropsAutoloadOut {
		return EnablePropsAutoloadOut{
			Properties: props,
		}
	})
}

type PropertiesLoaderIn struct {
	fx.In
	Properties []config.Properties `group:"properties"`
}

func NewPropertiesLoader(in PropertiesLoaderIn, options ...Option) (config.Loader, error) {
	option := new(config.Option)
	option.ActiveProfiles = utils.SliceFromCommaString(os.Getenv("APPLICATION_ENV"))
	option.ConfigPaths = utils.SliceFromCommaString(os.Getenv("APPLICATION_CONFIG_PATHS"))
	option.ConfigFormat = os.Getenv("APPLICATION_CONFIG_FORMAT")
	option.DebugFunc = log.Debugf
	for _, optFunc := range options {
		optFunc(option)
	}
	loader, err := config.NewLoader(*option, in.Properties)
	if err != nil {
		return nil, err
	}
	return loader, nil
}

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

func WithDebugLog(debugFunc config.DebugFunc) Option {
	return func(option *config.Option) {
		option.DebugFunc = debugFunc
	}
}
