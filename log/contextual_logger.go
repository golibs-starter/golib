package log

import (
	"context"
)

type ContextualLogger interface {
	// Debugc formats the message according to the format specifier
	// with some additional info in the context and logs it at debug level.
	Debugc(ctx context.Context, template string, args ...interface{})

	// Infoc formats the message according to the format specifier
	// with some additional info in the context and logs it at info level.
	Infoc(ctx context.Context, template string, args ...interface{})

	// Warnc formats the message according to the format specifier
	// with some additional info in the context and logs it at warn level.
	Warnc(ctx context.Context, template string, args ...interface{})

	// Errorc formats the message according to the format specifier
	// with some additional info in the context and logs it at error level.
	Errorc(ctx context.Context, template string, args ...interface{})

	// Fatalc formats the message according to the format specifier
	// with some additional info in the context and calls os.Exit.
	Fatalc(ctx context.Context, template string, args ...interface{})

	// Panicc formats the message according to the format specifier
	// with some additional info in the context and panics.
	Panicc(ctx context.Context, template string, args ...interface{})
}
