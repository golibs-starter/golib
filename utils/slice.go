package utils

import "strings"

func ContainsString(slice []string, needle string) bool {
	for _, n := range slice {
		if needle == n {
			return true
		}
	}
	return false
}

func PrependString(slice []string, e string) []string {
	slice = append(slice, "")
	copy(slice[1:], slice)
	slice[0] = e
	return slice
}

func SliceFromCommaString(str string) []string {
	slice := make([]string, 0)
	for _, env := range strings.Split(str, ",") {
		env = strings.TrimSpace(env)
		if env != "" {
			slice = append(slice, strings.TrimSpace(env))
		}
	}
	return slice
}
