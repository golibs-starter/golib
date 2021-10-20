package utils

import (
	"fmt"
	"reflect"
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

// MergeCaseInsensitiveMaps merges two maps
func MergeCaseInsensitiveMaps(src, tgt map[string]interface{}) map[string]interface{} {
	out := tgt
	for sk, sv := range src {
		tk := mapKeyExists(sk, tgt)
		if tk == "" {
			out[sk] = sv
			continue
		}
		tv, ok := tgt[tk]
		if !ok {
			out[sk] = sv
			continue
		}
		if msv, ok := sv.(map[interface{}]interface{}); ok {
			sv = castToMapStringInterface(msv)
		}
		if mtv, ok := tv.(map[interface{}]interface{}); ok {
			tv = castToMapStringInterface(mtv)
		}
		svType := reflect.TypeOf(sv)
		tvType := reflect.TypeOf(tv)
		if tvType != nil && svType != tvType {
			out[sk] = sv // want to override source value
			continue
		}

		switch ttv := tv.(type) {
		case map[string]interface{}:
			out[tk] = MergeCaseInsensitiveMaps(sv.(map[string]interface{}), ttv)
			break
		case []interface{}:
			for i, svItem := range sv.([]interface{}) {
				if msvItem, ok1 := svItem.(map[interface{}]interface{}); ok1 {
					svItem = castToMapStringInterface(msvItem)
				}
				if i >= len(ttv) {
					ttv = append(ttv, svItem)
					continue
				}
				if mTtvItem, ok1 := ttv[i].(map[interface{}]interface{}); ok1 {
					ttv[i] = castToMapStringInterface(mTtvItem)
				}
				svItemType := reflect.TypeOf(svItem)
				tvItemType := reflect.TypeOf(ttv[i])
				if tvItemType != nil && svItemType != tvItemType {
					ttv[i] = svItem
					continue
				}
				switch castedSvItem := svItem.(type) {
				case map[string]interface{}:
					ttv[i] = MergeCaseInsensitiveMaps(castedSvItem, ttv[i].(map[string]interface{}))
					break
				default:
					ttv[i] = svItem
				}
			}
			out[tk] = ttv
		default:
			out[tk] = sv
		}
	}
	return out
}

func castToMapStringInterface(src map[interface{}]interface{}) map[string]interface{} {
	tgt := map[string]interface{}{}
	for k, v := range src {
		tgt[fmt.Sprintf("%v", k)] = v
	}
	return tgt
}

func mapKeyExists(k string, m map[string]interface{}) string {
	lk := strings.ToLower(k)
	for mk := range m {
		lmk := strings.ToLower(mk)
		if lmk == lk {
			return mk
		}
	}
	return ""
}
