package log

import (
	"context"
	"gitlab.com/golibs-starter/golib/log/field"
)

type Logger interface {
	StdLogger

	WithCtx(ctx context.Context) Logger
	WithField(fields ...field.Field) Logger
}
