package config

type Properties interface {
	Prefix() string
}

type ApplicationProperties struct {
	Name string `mapstructure:"name" default:"unspecified"`
}

func (a ApplicationProperties) Prefix() string {
	return "vinid.application"
}
