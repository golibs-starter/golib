package log

const (
	OutputModeJson    = "json"
	OutputModeConsole = "console"
)

type FiledKey int

const (
	FieldKeyErr FiledKey = iota
	FieldKeyTime
	FieldKeyLevel
	FieldKeyName
	FieldKeyCaller
	FieldKeyFunc
	FieldKeyMsg
	FieldKeyStacktrace
)

var defaultFieldKeyMap = map[FiledKey]string{
	FieldKeyErr:        "error",
	FieldKeyTime:       "ts",
	FieldKeyLevel:      "level",
	FieldKeyName:       "logger",
	FieldKeyCaller:     "caller",
	FieldKeyFunc:       "", //Omit
	FieldKeyMsg:        "msg",
	FieldKeyStacktrace: "stacktrace",
}

type Options struct {
	// Development puts the logger in development mode, which changes the
	// behavior of DPanicLevel and takes stack traces more liberally.
	Development bool

	// LogLevel is the minimum enabled logging level.
	// In Development mode, LogLevel will be set to DEBUG,
	// Opposite, LogLevel will be set to INFO mode automatically.
	LogLevel string

	// JsonOutputMode Enable or disable json output mode.
	JsonOutputMode bool

	// DisableCaller stops annotating logs with the calling function's file
	// name and line number. By default, all logs are annotated.
	DisableCaller bool

	// DisableStacktrace completely disables automatic stacktrace capturing. By
	// default, stacktrace-s are captured for WarnLevel and above logs in
	// development and ErrorLevel and above in production.
	DisableStacktrace bool

	// CallerSkip Set the number of callers
	// will be skipped before show caller
	CallerSkip int

	// FieldKeyMap Set the keys used for each log entry.
	// If any key is empty, that portion of the entry is omitted.
	FieldKeyMap map[FiledKey]string

	// ContextExtractors Define the list of extractors
	// that will be used when extract value from log context.
	ContextExtractors ContextExtractors
}

type OptionFunc func(opt *Options)

func AddCallerSkip(skip int) OptionFunc {
	return func(opt *Options) {
		opt.CallerSkip = skip
	}
}
