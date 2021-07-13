package config

import "gitlab.id.vin/vincart/golib/utils"

const (
	DefaultProfile       = "default"
	DefaultActiveProfile = "local"
	DefaultConfigFormat  = "yaml"
	DefaultConfigPath    = "./config"
)

type Option struct {
	activeProfiles []string
	configPaths    []string
	configFormat   string // yaml, json
}

func setDefaultOption(option *Option) {
	if option.activeProfiles == nil {
		option.activeProfiles = make([]string, 0)
	}
	if !utils.ContainsString(option.activeProfiles, DefaultProfile) {
		option.activeProfiles = utils.PrependString(option.activeProfiles, DefaultProfile)
	}

	if !utils.ContainsString(option.activeProfiles, DefaultActiveProfile) {
		option.activeProfiles = append(option.activeProfiles, DefaultActiveProfile)
	}

	if len(option.configFormat) == 0 {
		option.configFormat = DefaultConfigFormat
	}

	if len(option.configPaths) == 0 {
		option.configPaths = []string{DefaultConfigPath}
	}
}
