package log

import (
	"context"
	"gitlab.com/golibs-starter/golib/log/field"
)

type Logger interface {
	StdLogger
	ContextualLogger

	WithCtx(ctx context.Context, additionalFields ...field.Field) Logger
	WithField(fields ...field.Field) Logger
	WithError(err error) Logger
	WithErrors(errs ...error) Logger
	WithAny(key string, value interface{}) Logger
}
