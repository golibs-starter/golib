package utils

import (
	"fmt"
	"os"
	"strings"
)

// ReplacePlaceholder Replaces a value in placeholder format
// by new value configured in environment variable.
//
// Placeholder format: ${EXAMPLE_VAR}
func ReplacePlaceholder(val interface{}) (interface{}, error) {
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
