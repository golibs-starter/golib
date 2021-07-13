package config

import (
	"fmt"
	"github.com/spf13/viper"
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

func NewLoader(option Option, debugLog func(msgFormat string, args ...interface{})) *ViperLoader {
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
		if err := l.viper.UnmarshalKey(properties.Prefix(), properties); err != nil {
			panic(fmt.Sprintf("[GoLib-error] Fatal error when binding config key [%s] to [%s]: %v",
				properties.Prefix(), propertiesName, err))
		}
		l.debugLog("[GoLib-debug] Properties [%s] loaded with prefix [%s]", propertiesName, properties.Prefix())
	}
}

func loadViper(option Option, debugLog func(msgFormat string, args ...interface{})) *viper.Viper {
	debugActiveProfiles := strings.Join(option.activeProfiles, ", ")
	debugPaths := strings.Join(option.configPaths, ", ")
	debugLog("[GoLib-debug] Loading active profiles [%s] in paths [%s] with format [%s]",
		debugActiveProfiles, debugPaths, option.configFormat)

	vi := viper.New()
	for _, activeProfile := range option.activeProfiles {
		vi.SetConfigName(activeProfile)
		vi.SetConfigType(option.configFormat)
		for _, path := range option.configPaths {
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
	return vi
}
