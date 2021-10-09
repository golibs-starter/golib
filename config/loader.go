package config

import (
	"bytes"
	"fmt"
	"github.com/creasty/defaults"
	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"reflect"
	"strings"
)

var keyDelimiter = "."

type Loader interface {
	Bind(properties ...Properties) error
}

type ViperLoader struct {
	viper            *viper.Viper
	option           Option
	properties       []Properties
	groupPropsConfig map[string]interface{}
}

func NewLoader(option Option, properties []Properties) (Loader, error) {
	setDefaultOption(&option)
	vi, err := loadViper(option, properties)
	if err != nil {
		return nil, err
	}
	return &ViperLoader{
		viper:            vi,
		option:           option,
		properties:       properties,
		groupPropsConfig: groupPropertiesValues(vi, properties),
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

		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			Result:           props,
			WeaklyTypedInput: true,
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				mapstructure.StringToTimeDurationHookFunc(),
				mapstructure.StringToSliceHookFunc(","),
				MapStructurePlaceholderValueHook(),
			),
		})
		if err != nil {
			return fmt.Errorf("[GoLib-error] Fatal error when init decoder for key [%s] to [%s]: %v",
				props.Prefix(), propsName, err)
		}
		if err := decoder.Decode(l.groupPropsConfig[props.Prefix()]); err != nil {
			return fmt.Errorf("[GoLib-error] Fatal error when binding config key [%s] to [%s]: %v",
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

func loadViper(option Option, propertiesList []Properties) (*viper.Viper, error) {
	option.DebugFunc("[GoLib-debug] Loading active profiles [%s] in paths [%s] with format [%s]",
		strings.Join(option.ActiveProfiles, ", "), strings.Join(option.ConfigPaths, ", "), option.ConfigFormat)

	vi := viper.NewWithOptions(viper.KeyDelimiter(keyDelimiter))
	vi.SetEnvKeyReplacer(strings.NewReplacer(keyDelimiter, "_"))
	vi.AutomaticEnv()

	if err := discoverDefaultValue(vi, propertiesList, option.DebugFunc); err != nil {
		return nil, err
	}

	if err := discoverActiveProfiles(vi, option); err != nil {
		return nil, err
	}
	return vi, nil
}

// discoverDefaultValue Discover default values for multiple properties at once
func discoverDefaultValue(vi *viper.Viper, propertiesList []Properties, debugFunc DebugFunc) error {
	for _, props := range propertiesList {
		propsName := reflect.TypeOf(props).String()

		// Set default value if its missing
		if err := defaults.Set(props); err != nil {
			return fmt.Errorf("[GoLib-error] Fatal error when set default values for [%s]: %v", propsName, err)
		}

		propsMap := structs.Map(props)
		defaultMap := wrapKeysAroundMap(strings.Split(props.Prefix(), keyDelimiter), propsMap, nil)

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
		debugFunc("[GoLib-debug] Default value was discovered for properties [%s]", propsName)
	}
	return nil
}

// discoverActiveProfiles Discover values for multiple active profiles at once
func discoverActiveProfiles(vi *viper.Viper, option Option) error {
	debugPaths := strings.Join(option.ConfigPaths, ", ")
	ext := option.ConfigFormat
	if ext == "yaml" {
		ext = "yml"
	}
	for _, path := range option.ConfigPaths {
		for _, activeProfile := range option.ActiveProfiles {
			filepath := path + "/" + activeProfile + "." + ext
			fileReader := NewFileReader(filepath, option.ConfigFormat, keyDelimiter)
			cfMap, err := fileReader.Read()
			if err != nil {
				return err
			}
			if err := vi.MergeConfigMap(cfMap); err != nil {
				return fmt.Errorf("[GoLib-error] Error when read active profile [%s] in paths [%s]: %v",
					activeProfile, debugPaths, err)
			}
			option.DebugFunc("[GoLib-debug] Active profile [%s] was loaded", activeProfile)
		}
	}
	return nil
}

func groupPropertiesValues(vi *viper.Viper, propertiesList []Properties) map[string]interface{} {
	allSettings := vi.AllSettings()
	group := make(map[string]interface{})
	for _, props := range propertiesList {
		m := deepSearchInMap(allSettings, props.Prefix())
		correctedVal, exists := correctSliceValues(vi, props.Prefix(), m)
		if exists {
			group[props.Prefix()] = correctedVal
		} else {
			group[props.Prefix()] = m
		}
	}
	return group
}

func correctSliceValues(vi *viper.Viper, prefix string, val interface{}) (interface{}, bool) {
	if slice, ok := val.([]interface{}); ok {
		for k, v := range slice {
			correctedVal, exists := correctSliceValues(vi, fmt.Sprintf("%s%s%d", prefix, keyDelimiter, k), v)
			if exists {
				slice[k] = correctedVal
			}
		}
	} else if m, ok := val.(map[interface{}]interface{}); ok {
		for k, v := range m {
			correctedVal, exists := correctSliceValues(vi, fmt.Sprintf("%s%s%s", prefix, keyDelimiter, k), v)
			if exists {
				m[k] = correctedVal
			}
		}
	} else if m, ok := val.(map[string]interface{}); ok {
		for k, v := range m {
			correctedVal, exists := correctSliceValues(vi, fmt.Sprintf("%s%s%s", prefix, keyDelimiter, k), v)
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
