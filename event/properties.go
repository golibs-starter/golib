package event

import "gitlab.id.vin/vincart/golib/config"

func NewProperties(loader config.Loader) (*Properties, error) {
	props := Properties{}
	err := loader.Bind(&props)
	return &props, err
}

type Properties struct {
	Log LogProperties
}

func (p Properties) Prefix() string {
	return "app.event"
}

type LogProperties struct {
	NotLogPayloadForEvents []string
}
