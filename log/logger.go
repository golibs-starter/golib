package log

import "context"

type Logger interface {
	StdLogger

	WithCtx(ctx context.Context) Logger
	WithField(fields ...Field) Logger
}
