package config

func NewAppProperties(loader Loader) (*AppProperties, error) {
	props := AppProperties{}
	err := loader.Bind(&props)
	return &props, err
}

type AppProperties struct {
	Name string `default:"unspecified"`
	Port int    `default:"8080"`
	Path string `default:"/"`
}

func (a AppProperties) Prefix() string {
	return "app"
}
