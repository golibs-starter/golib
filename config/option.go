package config

import "gitlab.id.vin/vincart/golib/utils"

const (
	DefaultProfile       = "default"
	DefaultActiveProfile = "local"
	DefaultConfigFormat  = "yaml"
	DefaultConfigPath    = "./config"
)

type Option struct {
	ActiveProfiles []string
	ConfigPaths    []string
	ConfigFormat   string // yaml, json
}

func setDefaultOption(option *Option) {
	if option.ActiveProfiles == nil {
		option.ActiveProfiles = make([]string, 0)
	}
	if !utils.ContainsString(option.ActiveProfiles, DefaultProfile) {
		option.ActiveProfiles = utils.PrependString(option.ActiveProfiles, DefaultProfile)
	}

	if !utils.ContainsString(option.ActiveProfiles, DefaultActiveProfile) {
		option.ActiveProfiles = append(option.ActiveProfiles, DefaultActiveProfile)
	}

	if len(option.ConfigFormat) == 0 {
		option.ConfigFormat = DefaultConfigFormat
	}

	if len(option.ConfigPaths) == 0 {
		option.ConfigPaths = []string{DefaultConfigPath}
	}
}
