package log

import (
	"gitlab.id.vin/vincart/golib/config"
)

func NewLoggingProperties(loader config.Loader) (*LoggingProperties, error) {
	props := LoggingProperties{}
	err := loader.Bind(&props)
	return &props, err
}

type LoggingProperties struct {
	Development    bool `default:"false"`
	JsonOutputMode bool `default:"true"`
	CallerSkip     int  `default:"4"`
}

func (l LoggingProperties) Prefix() string {
	return "app.logging"
}
