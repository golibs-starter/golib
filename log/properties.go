package log

import (
	"github.com/golibs-starter/golib/config"
)

func NewProperties(loader config.Loader) (*Properties, error) {
	props := Properties{}
	err := loader.Bind(&props)
	return &props, err
}

type Properties struct {
	// Development puts the logger in development mode, which changes the
	// behavior of DPanicLevel and takes stack traces more liberally.
	Development bool `default:"false"`

	// LogLevel is the minimum enabled logging level.
	// In Development mode, LogLevel will be set to DEBUG,
	// Opposite, LogLevel will be set to INFO mode automatically.
	LogLevel string

	// JsonOutputMode Enable or disable json output mode.
	JsonOutputMode bool `default:"true"`

	// DisableCaller stops annotating logs with the calling function's file
	// name and line number. By default, all logs are annotated.
	DisableCaller bool

	// DisableStacktrace completely disables automatic stacktrace capturing. By
	// default, stacktrace-s are captured for WarnLevel and above logs in
	// development and ErrorLevel and above in production.
	DisableStacktrace bool

	// CallerSkip Set the number of callers
	// will be skipped before show caller
	CallerSkip int `default:"1"`
}

func (l Properties) Prefix() string {
	return "app.logging"
}
