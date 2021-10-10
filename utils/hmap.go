package utils

import (
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"gopkg.in/yaml.v2"
	"strings"
)

type MapItem struct {
	Key   interface{}
	Value interface{}
}

func NewMapItem(key interface{}, value interface{}) MapItem {
	return MapItem{Key: key, Value: value}
}

func LinkedHMap(items ...MapItem) *linkedhashmap.Map {
	linkedMap := linkedhashmap.New()
	for _, item := range items {
		linkedMap.Put(item.Key, item.Value)
	}
	return linkedMap
}

func LinkedHMapToMapStr(hMap *linkedhashmap.Map) map[string]interface{} {
	mp := make(map[string]interface{})
	it := hMap.Iterator()
	for it.Next() {
		if k, ok1 := it.Key().(string); ok1 {
			v := it.Value()
			switch vOk := v.(type) {
			case *linkedhashmap.Map:
				mp[k] = LinkedHMapToMapStr(vOk)
				break
			case []*linkedhashmap.Map:
				m := make([]interface{}, 0)
				for _, hMp := range vOk {
					m = append(m, LinkedHMapToMapStr(hMp))
				}
				mp[k] = m
			default:
				mp[k] = it.Value()
			}
		}
	}
	return mp
}

func MergeLinkedHMap(dst *linkedhashmap.Map, src *linkedhashmap.Map) {
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
		MergeLinkedHMap(dstM, srcM)
		dst.Put(k, dstM)
	}
}

func ExpandInlineKeyInLinkedHMap(hMap *linkedhashmap.Map, delim string) *linkedhashmap.Map {
	var newSortedMp = linkedhashmap.New()
	for _, ki := range hMap.Keys() {
		k := ki.(string)
		v, _ := hMap.Get(k)
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
			expandedV := ExpandInlineKeyInLinkedHMap(vOk, delim)
			if existing, found := newSortedMp.Get(startK); found {
				if existingV, ok := existing.(*linkedhashmap.Map); ok {
					MergeLinkedHMap(existingV, expandedV)
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

func YamlMapSliceToLinkedHMap(ms yaml.MapSlice) *linkedhashmap.Map {
	m := linkedhashmap.New()
	for _, item := range ms {
		switch vOk := item.Value.(type) {
		case yaml.MapItem:
			m.Put(item.Key, YamlMapSliceToLinkedHMap(yaml.MapSlice{vOk}))
			break
		case yaml.MapSlice:
			m.Put(item.Key, YamlMapSliceToLinkedHMap(vOk))
			break
		case []yaml.MapSlice:
			slice := make([]*linkedhashmap.Map, 0)
			for _, child := range vOk {
				slice = append(slice, YamlMapSliceToLinkedHMap(child))
			}
			m.Put(item.Key, slice)
			break
		case []interface{}:
			sHashMap := make([]*linkedhashmap.Map, 0)
			sInf := make([]interface{}, 0)
			for _, child := range vOk {
				if childMs, ok := child.(yaml.MapSlice); ok {
					sHashMap = append(sHashMap, YamlMapSliceToLinkedHMap(childMs))
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
