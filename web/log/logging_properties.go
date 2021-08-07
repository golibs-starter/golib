package log

import "gitlab.id.vin/vincart/golib/config"

type LoggingProperties struct {
	Development    bool `default:"false"`
	JsonOutputMode bool `default:"true"`
	CallerSkip     int  `default:"2"`
}

func NewLoggingProperties(loader config.Loader) *LoggingProperties {
	props := LoggingProperties{}
	loader.Bind(&props)
	return &props
}

func (l LoggingProperties) Prefix() string {
	return "application.logging"
}
