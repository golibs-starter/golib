package config

import (
	"bytes"
	"fmt"
	"github.com/creasty/defaults"
	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gitlab.id.vin/vincart/golib/utils"
	"gopkg.in/yaml.v2"
	"reflect"
	"strings"
)

type Loader interface {
	Bind(properties ...Properties) error
}

type ViperLoader struct {
	viper          *viper.Viper
	option         Option
	groupedConfig  map[string]interface{}
	decodeHookFunc mapstructure.DecodeHookFunc
}

func NewLoader(option Option, properties []Properties) (Loader, error) {
	setDefaultOption(&option)
	reader, err := NewDefaultProfileReader(option.ConfigPaths, option.ConfigFormat, option.KeyDelimiter)
	if err != nil {
		return nil, err
	}
	vi, err := loadViper(reader, option, properties)
	if err != nil {
		return nil, err
	}
	return &ViperLoader{
		viper:         vi,
		option:        option,
		groupedConfig: groupPropertiesConfig(vi, properties, option),
		decodeHookFunc: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			MapStructurePlaceholderValueHook(),
		),
	}, nil
}

func (l *ViperLoader) Bind(propertiesList ...Properties) error {
	for _, props := range propertiesList {
		propsName := reflect.TypeOf(props).String()
		// Run pre-binding life cycle
		if propsPreBind, ok := props.(PropertiesPreBinding); ok {
			if err := propsPreBind.PreBinding(); err != nil {
				return err
			}
		}

		if err := l.decode(props); err != nil {
			return fmt.Errorf("[GoLib-error] Fatal error when decode config key [%s] to [%s]: %v",
				props.Prefix(), propsName, err)
		}

		// Run post-binding life cycle
		if propsPostBind, ok := props.(PropertiesPostBinding); ok {
			if err := propsPostBind.PostBinding(); err != nil {
				return err
			}
		}
		l.option.DebugFunc("[GoLib-debug] LoggingProperties [%s] loaded with prefix [%s]", propsName, props.Prefix())
	}
	return nil
}

func (l ViperLoader) decode(props Properties) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:       l.decodeHookFunc,
		WeaklyTypedInput: true,
		Result:           props,
	})
	if err != nil {
		return err
	}
	if err := decoder.Decode(l.groupedConfig[props.Prefix()]); err != nil {
		return err
	}
	return nil
}

func loadViper(reader ProfileReader, option Option, propertiesList []Properties) (*viper.Viper, error) {
	option.DebugFunc("[GoLib-debug] Loading active profiles [%s] in paths [%s] with format [%s]",
		strings.Join(option.ActiveProfiles, ", "), strings.Join(option.ConfigPaths, ", "), option.ConfigFormat)

	vi := viper.NewWithOptions(viper.KeyDelimiter(option.KeyDelimiter))
	vi.SetEnvKeyReplacer(strings.NewReplacer(option.KeyDelimiter, "_"))
	vi.AutomaticEnv()

	if err := discoverDefaultValue(vi, propertiesList, option); err != nil {
		return nil, err
	}

	if err := discoverActiveProfiles(vi, reader, option); err != nil {
		return nil, err
	}
	return vi, nil
}

// discoverDefaultValue Discover default values for multiple properties at once
func discoverDefaultValue(vi *viper.Viper, propertiesList []Properties, option Option) error {
	for _, props := range propertiesList {
		propsName := reflect.TypeOf(props).String()

		// Set default value if its missing
		if err := defaults.Set(props); err != nil {
			return fmt.Errorf("[GoLib-error] Fatal error when set default values for [%s]: %v", propsName, err)
		}

		propsMap := structs.Map(props)
		defaultMap := utils.WrapKeysAroundMap(strings.Split(props.Prefix(), option.KeyDelimiter), propsMap, nil)

		// set default values in viper.
		// Viper needs to know if a key exists in order to override it.
		// https://github.com/spf13/viper/issues/188
		b, err := yaml.Marshal(defaultMap)
		if err != nil {
			return err
		}
		vi.SetConfigType("yaml")
		if err := vi.MergeConfig(bytes.NewReader(b)); err != nil {
			return fmt.Errorf("[GoLib-error] Error when discover default value for properties [%s]: %v", propsName, err)
		}
		option.DebugFunc("[GoLib-debug] Default value was discovered for properties [%s]", propsName)
	}
	return nil
}

// discoverActiveProfiles Discover values for multiple active profiles at once
func discoverActiveProfiles(vi *viper.Viper, reader ProfileReader, option Option) error {
	debugPaths := strings.Join(option.ConfigPaths, ", ")
	for _, activeProfile := range option.ActiveProfiles {
		cfMap, err := reader.Read(activeProfile)
		if err != nil {
			return err
		}
		if err := vi.MergeConfigMap(cfMap); err != nil {
			return fmt.Errorf("[GoLib-error] Error when read active profile [%s] in paths [%s]: %v",
				activeProfile, debugPaths, err)
		}
		option.DebugFunc("[GoLib-debug] Active profile [%s] was loaded", activeProfile)
	}
	return nil
}

func groupPropertiesConfig(vi *viper.Viper, propertiesList []Properties, option Option) map[string]interface{} {
	allSettings := vi.AllSettings()
	group := make(map[string]interface{})
	for _, props := range propertiesList {
		m := utils.DeepSearchInMap(allSettings, props.Prefix(), option.KeyDelimiter)
		correctedVal, exists := correctSliceValues(vi, props.Prefix(), option.KeyDelimiter, m)
		if exists {
			group[props.Prefix()] = correctedVal
		} else {
			group[props.Prefix()] = m
		}
	}
	return group
}

func correctSliceValues(vi *viper.Viper, prefix string, delim string, val interface{}) (interface{}, bool) {
	if slice, ok := val.([]interface{}); ok {
		for k, v := range slice {
			correctedVal, exists := correctSliceValues(vi, fmt.Sprintf("%s%s%d", prefix, delim, k), delim, v)
			if exists {
				slice[k] = correctedVal
			}
		}
	} else if m, ok := val.(map[interface{}]interface{}); ok {
		for k, v := range m {
			correctedVal, exists := correctSliceValues(vi, fmt.Sprintf("%s%s%s", prefix, delim, k), delim, v)
			if exists {
				m[k] = correctedVal
			}
		}
	} else if m, ok := val.(map[string]interface{}); ok {
		for k, v := range m {
			correctedVal, exists := correctSliceValues(vi, fmt.Sprintf("%s%s%s", prefix, delim, k), delim, v)
			if exists {
				m[k] = correctedVal
			}
		}
	} else {
		if correctedVal := vi.Get(prefix); correctedVal != nil {
			return correctedVal, true
		}
	}
	return nil, false
}
