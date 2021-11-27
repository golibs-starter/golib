package event

import "gitlab.com/golibs-starter/golib/config"

func NewProperties(loader config.Loader) (*Properties, error) {
	props := Properties{}
	err := loader.Bind(&props)
	return &props, err
}

type Properties struct {
	Log LogProperties
}

func (p Properties) Prefix() string {
	return "vinid.event"
}

type LogProperties struct {
	NotLogPayloadForEvents []string
}
