package config

func NewAppProperties(loader Loader) (*AppProperties, error) {
	props := AppProperties{}
	err := loader.Bind(&props)
	return &props, err
}

type AppProperties struct {
	Name string `mapstructure:"name" default:"unspecified"`
	Port int    `mapstructure:"port" default:"8080"`
	Path string `mapstructure:"path" default:"/"`
}

func (a AppProperties) Prefix() string {
	return "application"
}
