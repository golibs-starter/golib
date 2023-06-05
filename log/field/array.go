package field

import "go.uber.org/zap"

// Strings constructs a field that carries a slice of strings.
func Strings(key string, ss []string) Field {
	return zap.Strings(key, ss)
}

// Errors constructs a field that carries a slice of errors.
func Errors(key string, errs []error) Field {
	return zap.Errors(key, errs)
}
