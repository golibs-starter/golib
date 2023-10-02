package log

import (
	"context"
	"github.com/golibs-starter/golib/log/field"
)

type Logger interface {
	StdLogger
	ContextualLogger

	// WithCtx adds additional info in the context and
	// additional fields to the logging context.
	WithCtx(ctx context.Context, additionalFields ...field.Field) Logger

	// WithField adds a variadic number of fields to the logging context.
	WithField(fields ...field.Field) Logger

	// WithError adds an error with FieldKeyErr field to the logging context.
	WithError(err error) Logger

	// WithErrors adds a field with FieldKeyErr field that carries a slice of errors.
	WithErrors(errs ...error) Logger

	// WithAny adds a key and an arbitrary value and chooses the best way to represent
	// them as a field.
	WithAny(key string, value interface{}) Logger
}
