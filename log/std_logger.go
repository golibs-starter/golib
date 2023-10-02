package log

type StdLogger interface {
	// Debug logs the provided arguments at debug level.
	// Spaces are added between arguments when neither is a string.
	Debug(args ...interface{})

	// Info logs the provided arguments at debug level.
	// Spaces are added between arguments when neither is a string.
	Info(args ...interface{})

	// Warn logs the provided arguments at debug level.
	// Spaces are added between arguments when neither is a string.
	Warn(args ...interface{})

	// Error logs the provided arguments at debug level.
	// Spaces are added between arguments when neither is a string.
	Error(args ...interface{})

	// Fatal constructs a message with the provided arguments and calls os.Exit.
	// Spaces are added between arguments when neither is a string.
	Fatal(args ...interface{})

	// Panic constructs a message with the provided arguments and panics.
	// Spaces are added between arguments when neither is a string.
	Panic(args ...interface{})

	// Debugf formats the message according to the format specifier
	// and logs it at debug level.
	Debugf(template string, args ...interface{})

	// Infof formats the message according to the format specifier
	// and logs it at info level.
	Infof(template string, args ...interface{})

	// Warnf formats the message according to the format specifier
	// and logs it at warn level.
	Warnf(template string, args ...interface{})

	// Errorf formats the message according to the format specifier
	// and logs it at error level.
	Errorf(template string, args ...interface{})

	// Fatalf formats the message according to the format specifier
	// and calls os.Exit.
	Fatalf(template string, args ...interface{})

	// Panicf formats the message according to the format specifier
	// and panics.
	Panicf(template string, args ...interface{})
}
