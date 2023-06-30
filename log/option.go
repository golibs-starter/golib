package log

const (
	OutputModeJson    = "json"
	OutputModeConsole = "console"
)

type FiledKey int

const (
	FieldKeyErr FiledKey = iota
)

var defaultFieldKeyMap = map[FiledKey]string{
	FieldKeyErr: "error",
}

type Options struct {
	// Development puts the logger in development mode, which changes the
	// behavior of DPanicLevel and takes stack traces more liberally.
	Development bool

	// Enable json output mode
	JsonOutputMode bool

	// Skip number of callers before show caller
	CallerSkip int

	FieldKeyMap map[FiledKey]string

	ContextExtractors ContextExtractors
}

type OptionFunc func(opt *Options)

func AddCallerSkip(skip int) OptionFunc {
	return func(opt *Options) {
		opt.CallerSkip = skip
	}
}
