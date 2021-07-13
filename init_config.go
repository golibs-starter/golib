package golib

import "gitlab.id.vin/vincart/golib/config"

func InitConfig(option config.Option, bindingProperties []config.Properties) {
	config.NewLoader(option, nil).Load(bindingProperties)
}
