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
