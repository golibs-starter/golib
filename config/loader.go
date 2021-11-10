package config

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/zenthangplus/defaults"
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
	vi, err := loadViper(reader, option)
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

		if err := l.decodeWithDefaults(props); err != nil {
			return fmt.Errorf("[GoLib-error] Fatal error when decode config key [%s] to [%s]: %v",
				props.Prefix(), propsName, err)
		}

		// Run post-binding life cycle
		if propsPostBind, ok := props.(PropertiesPostBinding); ok {
			if err := propsPostBind.PostBinding(); err != nil {
				return err
			}
		}
		l.option.DebugFunc("[GoLib-debug] Properties [%s] was loaded with prefix [%s]", propsName, props.Prefix())
	}
	return nil
}

func (l ViperLoader) decodeWithDefaults(props Properties) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:       l.decodeHookFunc,
		WeaklyTypedInput: true,
		Result:           props,
	})
	if err != nil {
		return err
	}
	loadedCfMap, ok := l.groupedConfig[props.Prefix()].(map[string]interface{})
	if !ok {
		return errors.New("loaded config inside prefix is not a map")
	}
	if err := decoder.Decode(loadedCfMap); err != nil {
		return errors.WithMessage(err, "cannot decode props")
	}

	// Set default value if its missing
	if err := defaults.Set(props); err != nil {
		return errors.WithMessage(err, "cannot set default")
	}

	// Set defaults is not enough here, because all zero values will be replaced with default value.
	// Example:
	//  When your config is `store.open (default:true)`,
	//  if you want to set `store.open=false`, after run above code, it will replace with `true` by default.
	//  => This behavior is not expected.
	//
	// The idea is convert properties to map again, then merge loaded config to that map.
	propsMap := make(map[string]interface{})
	if err := mapstructure.Decode(props, &propsMap); err != nil {
		return errors.WithMessage(err, "cannot decode props to map")
	}

	// Because mapstructure cannot decode struct inside slice, so convert to yaml and unmarshal again
	propsMapBytes, err := yaml.Marshal(propsMap)
	if err != nil {
		return errors.WithMessage(err, "cannot encode propsMap to yaml")
	}
	var defaultProps map[string]interface{}
	if err := yaml.Unmarshal(propsMapBytes, &defaultProps); err != nil {
		return errors.WithMessage(err, "cannot unmarshal propsMapBytes")
	}

	newPropsMap := utils.MergeCaseInsensitiveMaps(loadedCfMap, defaultProps)
	if err := decoder.Decode(newPropsMap); err != nil {
		return errors.New("cannot decode props again")
	}
	return nil
}

func loadViper(reader ProfileReader, option Option) (*viper.Viper, error) {
	option.DebugFunc("[GoLib-debug] Loading active profiles [%s] in paths [%s] with format [%s]",
		strings.Join(option.ActiveProfiles, ", "), strings.Join(option.ConfigPaths, ", "), option.ConfigFormat)

	vi := viper.NewWithOptions(viper.KeyDelimiter(option.KeyDelimiter))
	vi.SetEnvKeyReplacer(strings.NewReplacer(option.KeyDelimiter, "_"))
	vi.AutomaticEnv()

	if err := discoverActiveProfiles(vi, reader, option); err != nil {
		return nil, err
	}
	return vi, nil
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
