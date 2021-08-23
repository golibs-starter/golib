package config

import (
	"bytes"
	"fmt"
	"github.com/creasty/defaults"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"os"
	"reflect"
	"strings"
)

var keyDelimiter = "."

type Loader interface {
	Bind(properties ...Properties) error
}

type ViperLoader struct {
	viper      *viper.Viper
	option     Option
	properties []Properties
}

func NewLoader(option Option, properties []Properties) (Loader, error) {
	setDefaultOption(&option)
	vi, err := loadViper(option, properties)
	if err != nil {
		return nil, err
	}
	return &ViperLoader{
		viper:      vi,
		option:     option,
		properties: properties,
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

		// Unmarshal from config file
		if err := l.viper.UnmarshalKey(props.Prefix(), props); err != nil {
			return fmt.Errorf("[GoLib-error] Fatal error when binding config key [%s] to [%s]: %v",
				props.Prefix(), propsName, err)
		}

		// Set default value if its missing
		if err := l.setDefaults(propsName, props); err != nil {
			return err
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

func (l *ViperLoader) setDefaults(propertiesName string, properties Properties) error {
	if err := defaults.Set(properties); err != nil {
		return fmt.Errorf("[GoLib-error] Fatal error when set default values for [%s]: %v", propertiesName, err)
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

	// High priority for environment variable.
	// This is workaround solution because viper does not
	// treat env vars the same as other config
	// See https://github.com/spf13/viper/issues/188#issuecomment-399518663
	//
	// Notes: Currently vi.AllKeys() doesn't support key for array item, such as: foo.bar.0.username,
	// so environment variable cannot overwrite these values, replace placeholder also not working
	// (using PropertiesPostBinding to replace placeholder as a workaround solution).
	// TODO Improve it or wait for viper in next version
	for _, key := range vi.AllKeys() {
		val := vi.Get(key)
		if newVal, err := ReplacePlaceholderValue(val); err != nil {
			return nil, err
		} else {
			val = newVal
		}
		vi.Set(key, val)
	}
	return vi, nil
}

// discoverDefaultValue Discover default values for multiple properties at once
func discoverDefaultValue(vi *viper.Viper, propertiesList []Properties, debugFunc DebugFunc) error {
	for _, props := range propertiesList {
		propsName := reflect.TypeOf(props).String()

		// set default values in viper.
		// Viper needs to know if a key exists in order to override it.
		// https://github.com/spf13/viper/issues/188
		b, err := yaml.Marshal(convertSliceToNestedMap(strings.Split(props.Prefix(), keyDelimiter), props, nil))
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
	for _, activeProfile := range option.ActiveProfiles {
		vi.SetConfigName(activeProfile)
		vi.SetConfigType(option.ConfigFormat)
		for _, path := range option.ConfigPaths {
			vi.AddConfigPath(path)
		}
		if err := vi.MergeInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return fmt.Errorf("[GoLib-error] Error when read active profile [%s] in paths [%s]: %v",
					activeProfile, debugPaths, err)
			}
			return fmt.Errorf("[GoLib-debug] Config file not found when read active profile [%s] in paths [%s]",
				activeProfile, debugPaths)
		}
		option.DebugFunc("[GoLib-debug] Active profile [%s] was loaded", activeProfile)
	}
	return nil
}

func convertSliceToNestedMap(paths []string, endVal interface{}, inMap map[interface{}]interface{}) map[interface{}]interface{} {
	if inMap == nil {
		inMap = map[interface{}]interface{}{}
	}
	if len(paths) == 0 {
		return inMap
	}
	if len(paths) == 1 {
		inMap[paths[0]] = endVal
		return inMap
	}
	inMap[paths[0]] = convertSliceToNestedMap(paths[1:], endVal, map[interface{}]interface{}{})
	return inMap
}

// ReplacePlaceholderValue Replaces a value in placeholder format
// by new value configured in environment variable.
//
// Placeholder format: ${EXAMPLE_VAR}
func ReplacePlaceholderValue(val interface{}) (interface{}, error) {
	strVal, ok := val.(string)
	if !ok {
		return val, nil
	}
	// Make sure the value starts with ${ and end with }
	if !strings.HasPrefix(strVal, "${") || !strings.HasSuffix(strVal, "}") {
		return val, nil
	}
	key := strings.TrimSuffix(strings.TrimPrefix(strVal, "${"), "}")
	if len(key) == 0 {
		return nil, fmt.Errorf("invalid config placeholder format. Expected ${EX_ENV}, got [%s]", strVal)
	}
	res, present := os.LookupEnv(key)
	if !present {
		return nil, fmt.Errorf("mandatory env variable not found [%s]", key)
	}
	return res, nil
}
