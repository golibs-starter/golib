package utils

import (
	"encoding/json"
)

func StructToMap(structVal interface{}) (map[string]interface{}, error) {
	var mapVal map[string]interface{}
	bytes, err := json.Marshal(structVal)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(bytes, &mapVal); err != nil {
		return nil, err
	}
	return mapVal, nil
}

func FlatStruct(structVal interface{}) ([]interface{}, error) {
	mapVal, err := StructToMap(structVal)
	if err != nil {
		return nil, err
	}
	keysAndValues := make([]interface{}, 0)
	for field, val := range mapVal {
		keysAndValues = append(keysAndValues, field, val)
	}
	return keysAndValues, nil
}
