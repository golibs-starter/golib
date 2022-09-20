package utils

import "reflect"

func GetStructShortName(val interface{}) string {
	if val == nil {
		return ""
	}
	if t := reflect.TypeOf(val); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else if t.Kind() == reflect.Struct {
		return t.Name()
	} else {
		return ""
	}
}

func GetStructFullname(val interface{}) string {
	if val == nil {
		return ""
	}
	if t := reflect.TypeOf(val); t.Kind() == reflect.Ptr {
		return t.Elem().String()
	} else if t.Kind() == reflect.Struct {
		return t.String()
	} else {
		return ""
	}
}
