package config

import (
	"fmt"
	"github.com/spf13/viper"
	"reflect"
	"strings"
)

type Loader struct {
	option   Option
	debugLog func(msgFormat string, args ...interface{})
}

func NewLoader(option Option, debugLog func(msgFormat string, args ...interface{})) *Loader {
	setDefaultOption(&option)
	if debugLog == nil {
		debugLog = func(msgFormat string, args ...interface{}) {
			_, _ = fmt.Printf(msgFormat+"\n", args...)
		}
	}
	return &Loader{
		option:   option,
		debugLog: debugLog,
	}
}

func (l Loader) Load(bindingProperties []Properties) {
	debugActiveProfiles := strings.Join(l.option.activeProfiles, ", ")
	debugPaths := strings.Join(l.option.configPaths, ", ")
	l.debugLog("[GoLib-debug] Loading active profiles [%s] in paths [%s] with format [%s]",
		debugActiveProfiles, debugPaths, l.option.configFormat)

	vi := viper.New()
	for _, activeProfile := range l.option.activeProfiles {
		vi.SetConfigName(activeProfile)
		vi.SetConfigType(l.option.configFormat)
		for _, path := range l.option.configPaths {
			vi.AddConfigPath(path)
		}
		if err := vi.MergeInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				l.debugLog("[GoLib-debug] Config file not found when read active profile [%s] in paths [%s]",
					activeProfile, debugPaths)
			} else {
				panic(fmt.Sprintf("[GoLib-error] Fatal error when read active profile [%s] in paths [%s]: %v",
					activeProfile, debugPaths, err))
			}
		} else {
			l.debugLog("[GoLib-debug] Active profile [%s] was loaded", activeProfile)
		}
	}
	l.bindProperties(vi, bindingProperties)
}

func (l *Loader) bindProperties(vi *viper.Viper, bindingProperties []Properties) {
	for _, properties := range bindingProperties {
		propertiesName := reflect.TypeOf(properties).String()
		if err := vi.UnmarshalKey(properties.Prefix(), properties); err != nil {
			panic(fmt.Sprintf("[GoLib-error] Fatal error when binding config key [%s] to [%s]: %v",
				properties.Prefix(), propertiesName, err))
		} else {
			l.debugLog("[GoLib-debug] Properties [%s] loaded with prefix [%s]",
				propertiesName, properties.Prefix())
		}
	}
}
