package field

import (
	"go.uber.org/zap"
	"time"
)

// Bools constructs a field that carries a slice of bools.
func Bools(key string, bs []bool) Field {
	return zap.Bools(key, bs)
}

// ByteStrings constructs a field that carries a slice of []byte, each of which
// must be UTF-8 encoded text.
func ByteStrings(key string, bss [][]byte) Field {
	return zap.ByteStrings(key, bss)
}

// Complex128s constructs a field that carries a slice of complex numbers.
func Complex128s(key string, nums []complex128) Field {
	return zap.Complex128s(key, nums)
}

// Complex64s constructs a field that carries a slice of complex numbers.
func Complex64s(key string, nums []complex64) Field {
	return zap.Complex64s(key, nums)
}

// Durations constructs a field that carries a slice of time.Durations.
func Durations(key string, ds []time.Duration) Field {
	return zap.Durations(key, ds)
}

// Float64s constructs a field that carries a slice of floats.
func Float64s(key string, nums []float64) Field {
	return zap.Float64s(key, nums)
}

// Float32s constructs a field that carries a slice of floats.
func Float32s(key string, nums []float32) Field {
	return zap.Float32s(key, nums)
}

// Ints constructs a field that carries a slice of integers.
func Ints(key string, nums []int) Field {
	return zap.Ints(key, nums)
}

// Int64s constructs a field that carries a slice of integers.
func Int64s(key string, nums []int64) Field {
	return zap.Int64s(key, nums)
}

// Int32s constructs a field that carries a slice of integers.
func Int32s(key string, nums []int32) Field {
	return zap.Int32s(key, nums)
}

// Int16s constructs a field that carries a slice of integers.
func Int16s(key string, nums []int16) Field {
	return zap.Int16s(key, nums)
}

// Int8s constructs a field that carries a slice of integers.
func Int8s(key string, nums []int8) Field {
	return zap.Int8s(key, nums)
}

// Strings constructs a field that carries a slice of strings.
func Strings(key string, ss []string) Field {
	return zap.Strings(key, ss)
}

// Times constructs a field that carries a slice of time.Times.
func Times(key string, ts []time.Time) Field {
	return zap.Times(key, ts)
}

// Uints constructs a field that carries a slice of unsigned integers.
func Uints(key string, nums []uint) Field {
	return zap.Uints(key, nums)
}

// Uint64s constructs a field that carries a slice of unsigned integers.
func Uint64s(key string, nums []uint64) Field {
	return zap.Uint64s(key, nums)
}

// Uint32s constructs a field that carries a slice of unsigned integers.
func Uint32s(key string, nums []uint32) Field {
	return zap.Uint32s(key, nums)
}

// Uint16s constructs a field that carries a slice of unsigned integers.
func Uint16s(key string, nums []uint16) Field {
	return zap.Uint16s(key, nums)
}

// Uint8s constructs a field that carries a slice of unsigned integers.
func Uint8s(key string, nums []uint8) Field {
	return zap.Uint8s(key, nums)
}

// Uintptrs constructs a field that carries a slice of pointer addresses.
func Uintptrs(key string, us []uintptr) Field {
	return zap.Uintptrs(key, us)
}

// Errors constructs a field that carries a slice of errors.
func Errors(key string, errs []error) Field {
	return zap.Errors(key, errs)
}
