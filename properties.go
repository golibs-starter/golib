package golib

import (
	"fmt"
	"github.com/golibs-starter/golib/config"
	"github.com/golibs-starter/golib/utils"
	"go.uber.org/fx"
	"log"
	"os"
	"reflect"
	"runtime"
	"strings"
)

func ProvideProps(propConstructor interface{}) fx.Option {
	sampleProps, err := makeSampleProperties(propConstructor)
	if err != nil {
		return fx.Error(err)
	}
	return fx.Options(
		fx.Provide(fx.Annotated{
			Group: "properties",
			Target: func() config.Properties {
				return sampleProps
			},
		}),
		fx.Provide(propConstructor),
	)
}

func ProvidePropsOption(option Option) fx.Option {
	return fx.Provide(fx.Annotated{
		Group: "properties_option",
		Target: func() Option {
			return option
		},
	})
}

func PropertiesOpt() fx.Option {
	return fx.Provide(NewPropertiesLoader)
}

type PropertiesLoaderIn struct {
	fx.In
	Properties []config.Properties `group:"properties"`
	Options    []Option            `group:"properties_option"`
}

func NewPropertiesLoader(in PropertiesLoaderIn) (config.Loader, error) {
	profiles := strings.TrimSpace(os.Getenv("APP_PROFILES"))
	if len(profiles) == 0 {
		profiles = strings.TrimSpace(os.Getenv("APP_ENV"))
	}

	// Set default option
	option := new(config.Option)
	option.ActiveProfiles = utils.SliceFromCommaString(profiles)
	option.ConfigPaths = utils.SliceFromCommaString(os.Getenv("APP_CONFIG_PATHS"))
	option.ConfigFormat = os.Getenv("APP_CONFIG_FORMAT")
	option.DebugFunc = log.Printf

	// Apply user option
	for _, optFunc := range in.Options {
		optFunc(option)
	}

	if len(option.ActiveProfiles) == 0 {
		option.ActiveProfiles = []string{"local"}
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

func makeSampleProperties(f interface{}) (config.Properties, error) {
	t := reflect.TypeOf(f)
	if t.Kind() != reflect.Func {
		return nil, fmt.Errorf("properties constructor must be a function. %s is provided", t.Name())
	}
	for i := 0; i < t.NumOut(); i++ {
		ele := t.Out(i)
		var val reflect.Value
		if ele.Kind() == reflect.Ptr {
			val = reflect.Zero(ele.Elem())
		} else {
			val = reflect.Zero(ele)
		}
		p := reflect.New(val.Type())
		p.Elem().Set(val)
		if props, ok := p.Interface().(config.Properties); ok {
			return props, nil
		}
	}
	return nil, fmt.Errorf("no properties found in output of constructor [%s]",
		runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name())
}
