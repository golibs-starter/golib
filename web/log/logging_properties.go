package log

import "gitlab.id.vin/vincart/golib/config"

type LoggingProperties struct {
	Development    bool `default:"false"`
	JsonOutputMode bool `default:"true"`
	CallerSkip     int  `default:"2"`
}

func NewLoggingProperties(loader config.Loader) (*LoggingProperties, error) {
	props := LoggingProperties{}
	if err := loader.Bind(&props); err != nil {
		return nil, err
	}
	return &props, nil
}

func (l LoggingProperties) Prefix() string {
	return "application.logging"
}
