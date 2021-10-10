package utils

import (
	"strings"
)

func DeepSearchInMap(m map[string]interface{}, key string, keyDelimiter string) map[string]interface{} {
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

func WrapKeysAroundMap(keyPaths []string, endVal interface{}, inMap map[string]interface{}) map[string]interface{} {
	if inMap == nil {
		inMap = map[string]interface{}{}
	}
	if len(keyPaths) == 0 {
		return inMap
	}
	if len(keyPaths) == 1 {
		inMap[keyPaths[0]] = endVal
		return inMap
	}
	inMap[keyPaths[0]] = WrapKeysAroundMap(keyPaths[1:], endVal, map[string]interface{}{})
	return inMap
}
