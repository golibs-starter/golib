package utils

import "github.com/emirpasic/gods/maps/linkedhashmap"

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
