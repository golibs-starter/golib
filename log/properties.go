package log

import (
	"gitlab.com/golibs-starter/golib/config"
)

func NewProperties(loader config.Loader) (*Properties, error) {
	props := Properties{}
	err := loader.Bind(&props)
	return &props, err
}

type Properties struct {
	Development    bool `default:"false"`
	JsonOutputMode bool `default:"true"`
	CallerSkip     int  `default:"1"`
}

func (l Properties) Prefix() string {
	return "app.logging"
}
