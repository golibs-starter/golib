package example

import (
	"gitlab.id.vin/vincart/golib/config"
	"time"
)

// ==================================================
// ======== Example about declare properties ========
// ==================================================

func NewSampleProperties(loader config.Loader) (*SampleProperties, error) {
	props := SampleProperties{}
	err := loader.Bind(&props)
	return &props, err
}

type SampleProperties struct {
	// We use github.com/go-playground/validator to validate properties
	Field1 string `validate:"required"`

	// We use https://github.com/zenthangplus/defaults to set default for properties
	Field2 int `default:"10"`

	// We use github.com/mitchellh/mapstructure to bind config to properties
	Field3 []time.Duration `mapstructure:"field3_new_name"`
}

// Prefix Defines the properties prefix
func (s SampleProperties) Prefix() string {
	return "app.sample"
}
