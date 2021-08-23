package config

import (
	"fmt"
	"gitlab.id.vin/vincart/golib/utils"
)

const (
	DefaultProfile      = "default"
	DefaultConfigFormat = "yaml"
	DefaultConfigPath   = "./config"
)

type DebugFunc func(msgFormat string, args ...interface{})

type Option struct {
	ActiveProfiles []string
	ConfigPaths    []string
	ConfigFormat   string // yaml, json
	DebugFunc      DebugFunc
}

func setDefaultOption(option *Option) {
	if option.ActiveProfiles == nil {
		option.ActiveProfiles = make([]string, 0)
	}
	if !utils.ContainsString(option.ActiveProfiles, DefaultProfile) {
		option.ActiveProfiles = utils.PrependString(option.ActiveProfiles, DefaultProfile)
	}

	if len(option.ConfigFormat) == 0 {
		option.ConfigFormat = DefaultConfigFormat
	}

	if len(option.ConfigPaths) == 0 {
		option.ConfigPaths = []string{DefaultConfigPath}
	}

	if option.DebugFunc == nil {
		option.DebugFunc = func(msgFormat string, args ...interface{}) {
			_, _ = fmt.Printf(msgFormat+"\n", args...)
		}
	}
}
