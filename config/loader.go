package config

import (
	"fmt"
	"github.com/creasty/defaults"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"strings"
)

type Loader interface {
	Bind(properties ...Properties)
}

type ViperLoader struct {
	option   Option
	debugLog func(msgFormat string, args ...interface{})
	viper    *viper.Viper
}

func NewLoader(option Option, debugLog func(msgFormat string, args ...interface{})) Loader {
	setDefaultOption(&option)
	if debugLog == nil {
		debugLog = func(msgFormat string, args ...interface{}) {
			_, _ = fmt.Printf(msgFormat+"\n", args...)
		}
	}
	return &ViperLoader{
		option:   option,
		debugLog: debugLog,
		viper:    loadViper(option, debugLog),
	}
}

func (l *ViperLoader) Bind(propertiesList ...Properties) {
	for _, properties := range propertiesList {
		propertiesName := reflect.TypeOf(properties).String()
		// Run pre-binding life cycle
		if propsPreBind, ok := properties.(PropertiesPreBinding); ok {
			propsPreBind.PreBinding()
		}

		// Unmarshal from config file
		if err := l.viper.UnmarshalKey(properties.Prefix(), properties); err != nil {
			panic(fmt.Sprintf("[GoLib-error] Fatal error when binding config key [%s] to [%s]: %v",
				properties.Prefix(), propertiesName, err))
		}

		// Set default value if its missing
		l.setDefaults(propertiesName, properties)

		// Run post-binding life cycle
		if propsPostBind, ok := properties.(PropertiesPostBinding); ok {
			propsPostBind.PostBinding()
		}
		l.debugLog("[GoLib-debug] Properties [%s] loaded with prefix [%s]", propertiesName, properties.Prefix())
	}
}

func (l *ViperLoader) setDefaults(propertiesName string, properties Properties) {
	if err := defaults.Set(properties); err != nil {
		panic(fmt.Sprintf("[GoLib-error] Fatal error when set default values for [%s]: %v", propertiesName, err))
	}
}

func loadViper(option Option, debugLog func(msgFormat string, args ...interface{})) *viper.Viper {
	debugActiveProfiles := strings.Join(option.ActiveProfiles, ", ")
	debugPaths := strings.Join(option.ConfigPaths, ", ")
	debugLog("[GoLib-debug] Loading active profiles [%s] in paths [%s] with format [%s]",
		debugActiveProfiles, debugPaths, option.ConfigFormat)

	vi := viper.New()
	vi.AutomaticEnv()
	vi.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	for _, activeProfile := range option.ActiveProfiles {
		vi.SetConfigName(activeProfile)
		vi.SetConfigType(option.ConfigFormat)
		for _, path := range option.ConfigPaths {
			vi.AddConfigPath(path)
		}
		if err := vi.MergeInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				panic(fmt.Sprintf("[GoLib-error] Fatal error when read active profile [%s] in paths [%s]: %v",
					activeProfile, debugPaths, err))
			}
			debugLog("[GoLib-debug] Config file not found when read active profile [%s] in paths [%s]",
				activeProfile, debugPaths)
			continue
		}
		debugLog("[GoLib-debug] Active profile [%s] was loaded", activeProfile)
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
			panic(err)
		} else {
			val = newVal
		}
		vi.Set(key, val)
	}
	return vi
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
