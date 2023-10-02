package log

import (
	"context"
)

type ContextualLogger interface {
	Debugc(ctx context.Context, template string, args ...interface{})
	Infoc(ctx context.Context, template string, args ...interface{})
	Warnc(ctx context.Context, template string, args ...interface{})
	Errorc(ctx context.Context, template string, args ...interface{})
	Fatalc(ctx context.Context, template string, args ...interface{})
	Panicc(ctx context.Context, template string, args ...interface{})
}
