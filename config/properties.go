package config

type Properties interface {
	Prefix() string
}

type PropertiesPreBinding interface {
	PreBinding() error
}

type PropertiesPostBinding interface {
	PostBinding() error
}

type AppProperties struct {
	Name string `mapstructure:"name" default:"unspecified"`
	Port int    `mapstructure:"port" default:"8080"`
	Path string `mapstructure:"path" default:"/"`
}

func NewApplicationProperties(loader Loader) (*AppProperties, error) {
	props := AppProperties{}
	if err := loader.Bind(&props); err != nil {
		return nil, err
	}
	return &props, nil
}

func (a AppProperties) Prefix() string {
	return "application"
}
