package config

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"os"
	"reflect"
	"strings"
)

func MapStructurePlaceholderValueHook() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		return ReplacePlaceholderValue(data.(string))
	}
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

func deepSearchInMap(m map[string]interface{}, key string) map[string]interface{} {
	parts := strings.Split(key, keyDelimiter)
	for _, part := range parts {
		val, ok := m[part]
		if !ok {
			return make(map[string]interface{})
		}
		m, ok = val.(map[string]interface{})
		if !ok {
			return make(map[string]interface{})
		}
	}
	return m
}

func wrapKeysAroundMap(paths []string, endVal interface{}, inMap map[string]interface{}) map[string]interface{} {
	if inMap == nil {
		inMap = map[string]interface{}{}
	}
	if len(paths) == 0 {
		return inMap
	}
	if len(paths) == 1 {
		inMap[paths[0]] = endVal
		return inMap
	}
	inMap[paths[0]] = wrapKeysAroundMap(paths[1:], endVal, map[string]interface{}{})
	return inMap
}

func mapToLowerKey(mp map[string]interface{}) map[string]interface{} {
	if mp == nil {
		return nil
	}
	newMap := make(map[string]interface{})
	for k, v := range mp {
		//lowerK := strings.ToLower(k)
		lowerK := k
		switch vOk := v.(type) {
		case map[string]interface{}:
			newMap[lowerK] = mapToLowerKey(vOk)
		default:
			newMap[lowerK] = v
		}
	}
	return newMap
}
