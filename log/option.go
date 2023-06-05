package log

import (
	"context"
	"gitlab.com/golibs-starter/golib/log/field"
)

const (
	OutputModeJson    = "json"
	OutputModeConsole = "console"
)

type Options struct {
	// Development puts the logger in development mode, which changes the
	// behavior of DPanicLevel and takes stack traces more liberally.
	Development bool

	// Enable json output mode
	JsonOutputMode bool

	// Skip number of callers before show caller
	CallerSkip int

	ContextExtractor func(ctx context.Context) []field.Field
}

type OptionFunc func(opt *Options)

func AddCallerSkip(skip int) OptionFunc {
	return func(opt *Options) {
		opt.CallerSkip = skip
	}
}
