package log

type LoggingProperties struct {
	Development    bool `mapstructure:"development"`
	JsonOutputMode bool `mapstructure:"json_output_mode"`
	CallerSkip     int  `mapstructure:"caller_skip" default:"2"`
}

func (l LoggingProperties) Prefix() string {
	return "vinid.logging"
}
