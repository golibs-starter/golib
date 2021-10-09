package config

import (
	"fmt"
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
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

func expandInlineKeyInMap_old(mp map[string]interface{}, delim string) map[string]interface{} {
	var newMp = make(map[string]interface{})
	for k, v := range mp {
		var startK string
		var startV interface{}
		hash := strings.Split(k, delim)
		if len(hash) == 1 {
			startK = hash[0]
			startV = v
		} else {
			startK = hash[0]
			startV = map[string]interface{}{
				strings.Join(hash[1:], delim): v,
			}
		}
		switch vOk := startV.(type) {
		case map[string]interface{}:
			newMp[startK] = expandInlineKeyInMap_old(vOk, delim)
		default:
			newMp[startK] = startV
		}
	}
	return newMp
}

func expandInlineKeyInMap_override1(mp map[string]interface{}, delim string) map[string]interface{} {
	var newMp = make(map[string]interface{})
	for k, v := range mp {
		var startK string
		var startV interface{}
		hash := strings.Split(k, delim)
		if len(hash) == 1 {
			startK = hash[0]
			startV = v
		} else {
			startK = hash[0]
			startV = map[string]interface{}{
				strings.Join(hash[1:], delim): v,
			}
		}
		switch vOk := startV.(type) {
		case map[string]interface{}:
			expandedV := expandInlineKeyInMap_override1(vOk, delim)
			if mergedV, ok := newMp[startK]; ok {
				if mergedMapV, isMap := mergedV.(map[string]interface{}); isMap {
					_ = mergo.Merge(&mergedMapV, expandedV, mergo.WithOverride)
					newMp[startK] = mergedMapV
				} else {
					newMp[startK] = expandedV
				}
			} else {
				newMp[startK] = expandedV
			}
		default:
			newMp[startK] = startV
		}
	}
	return newMp
}

func expandInlineKeyInMap(sortedMp *linkedhashmap.Map, delim string) *linkedhashmap.Map {
	var newSortedMp = linkedhashmap.New()
	for _, ki := range sortedMp.Keys() {
		k := ki.(string)
		v, _ := sortedMp.Get(k)
		var startK string
		var startV interface{}
		hash := strings.Split(k, delim)
		if len(hash) == 1 {
			startK = hash[0]
			startV = v
		} else {
			startK = hash[0]
			startM := linkedhashmap.New()
			startM.Put(strings.Join(hash[1:], delim), v)
			startV = startM
		}
		switch vOk := startV.(type) {
		case *linkedhashmap.Map:
			expandedV := expandInlineKeyInMap(vOk, delim)
			if existing, found := newSortedMp.Get(startK); found {
				if existingV, ok := existing.(*linkedhashmap.Map); ok {
					mergeLinkedHMap(existingV, expandedV)
				} else {
					newSortedMp.Put(startK, expandedV)
				}
			} else {
				newSortedMp.Put(startK, expandedV)
			}
		default:
			newSortedMp.Put(startK, vOk)
		}
	}
	return newSortedMp
}

func mergeLinkedHMap(dst *linkedhashmap.Map, src *linkedhashmap.Map) {
	for _, k := range src.Keys() {
		srcV, _ := src.Get(k)
		dstV, found := dst.Get(k)
		if !found {
			dst.Put(k, srcV)
			continue
		}
		srcM, ok1 := srcV.(*linkedhashmap.Map)
		dstM, ok2 := dstV.(*linkedhashmap.Map)
		if !ok1 || !ok2 {
			dst.Put(k, srcV)
			continue
		}
		mergeLinkedHMap(dstM, srcM)
		dst.Put(k, dstM)
	}
}

func yamlMapSliceToLinkedHMap(ms yaml.MapSlice) *linkedhashmap.Map {
	m := linkedhashmap.New()
	for _, item := range ms {
		switch vOk := item.Value.(type) {
		case yaml.MapItem:
			m.Put(item.Key, yamlMapSliceToLinkedHMap(yaml.MapSlice{vOk}))
			break
		case yaml.MapSlice:
			m.Put(item.Key, yamlMapSliceToLinkedHMap(vOk))
			break
		case []yaml.MapSlice:
			slice := make([]*linkedhashmap.Map, 0)
			for _, child := range vOk {
				slice = append(slice, yamlMapSliceToLinkedHMap(child))
			}
			m.Put(item.Key, slice)
			break
		case []interface{}:
			sHashMap := make([]*linkedhashmap.Map, 0)
			sInf := make([]interface{}, 0)
			for _, child := range vOk {
				if childMs, ok := child.(yaml.MapSlice); ok {
					sHashMap = append(sHashMap, yamlMapSliceToLinkedHMap(childMs))
				} else {
					sInf = append(sInf, child)
				}
			}
			if len(sHashMap) > 0 {
				m.Put(item.Key, sHashMap)
			} else {
				m.Put(item.Key, sInf)
			}
			break
		default:
			m.Put(item.Key, item.Value)
		}
	}
	return m
}
