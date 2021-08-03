package log

type LoggingProperties struct {
	Development    bool `default:"false"`
	JsonOutputMode bool `default:"true"`
	CallerSkip     int  `default:"2"`
}

func (l LoggingProperties) Prefix() string {
	return "application.logging"
}
