package config

type Properties interface {
	Prefix() string
}

type PropertiesPreBinding interface {
	PreBinding()
}

type PropertiesPostBinding interface {
	PostBinding()
}

type AppProperties struct {
	Name string `mapstructure:"name" default:"unspecified"`
	Port int    `mapstructure:"port" default:"8080"`
	Path string `mapstructure:"path" default:"/"`
}

func NewApplicationProperties(loader Loader) *AppProperties {
	props := AppProperties{}
	loader.Bind(&props)
	return &props
}

func (a AppProperties) Prefix() string {
	return "application"
}
