package config

import (
	"github.com/golibs-starter/golib/utils"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

func MapStructurePlaceholderValueHook() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		return utils.ReplacePlaceholder(data.(string))
	}
}
