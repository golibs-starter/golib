package event

import "github.com/golibs-starter/golib/config"

func NewProperties(loader config.Loader) (*Properties, error) {
	props := Properties{}
	err := loader.Bind(&props)
	return &props, err
}

type Properties struct {
	ChannelSize int `default:"10"`
	Log         LogProperties
}

func (p Properties) Prefix() string {
	return "app.event"
}

type LogProperties struct {
	NotLogPayloadForEvents []string
}
