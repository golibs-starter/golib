package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	assert "github.com/stretchr/testify/require"
	"testing"
)

type testStoreWithValidation struct {
	Name           string                      `validate:"required"`
	Path           string                      `default:"Apple/Central" validate:"required"`
	Tags           []string                    `validate:"len=2"`
	Products       []testProductWithValidation `validate:"min=1"`
	ProductCodeMap map[string]testProductWithValidation
}

func (t testStoreWithValidation) Prefix() string {
	return "org.storeWithValidation"
}

type testProductWithValidation struct {
	Code  string `validate:"required"`
	Title string `validate:"required"`
}

func TestLoaderBindingValidate_WhenConfigIsValid_ShouldReturnSuccess(t *testing.T) {
	loader, err := NewLoader(Option{
		ActiveProfiles: []string{"test_validation_success"},
		ConfigPaths:    []string{"./test_assets"},
		ConfigFormat:   "yml",
	}, []Properties{new(testStoreWithValidation)})
	assert.NoError(t, err)

	props := testStoreWithValidation{}
	err = loader.Bind(&props)
	assert.NoError(t, err)
	assert.Equal(t, "Apple Inc.", props.Name)
	assert.Equal(t, []string{"Iphone", "Ipad"}, props.Tags)
	assert.Len(t, props.Products, 2)
	assert.Equal(t, "IPHONE_6", props.Products[0].Code)
	assert.Equal(t, "Iphone 6", props.Products[0].Title)
	assert.Equal(t, "IPAD_MINI", props.Products[1].Code)
	assert.Equal(t, "Ipad mini", props.Products[1].Title)

	assert.Len(t, props.ProductCodeMap, 2)
	assert.Equal(t, "IPHONE_6", props.ProductCodeMap["iphone_6"].Code)
	assert.Equal(t, "Iphone 6", props.ProductCodeMap["iphone_6"].Title)
	assert.Equal(t, "IPAD_MINI", props.ProductCodeMap["ipad_mini"].Code)
	assert.Equal(t, "Ipad mini", props.ProductCodeMap["ipad_mini"].Title)
}

func TestLoaderBindingValidate_WhenConfigIsMissingField_ShouldReturnError(t *testing.T) {
	loader, err := NewLoader(Option{
		ActiveProfiles: []string{"test_validation_missing_field"},
		ConfigPaths:    []string{"./test_assets"},
		ConfigFormat:   "yml",
	}, []Properties{new(testStoreWithValidation)})
	assert.NoError(t, err)

	props := testStoreWithValidation{}
	err = loader.Bind(&props)
	assert.Error(t, err)
	cause := errors.Cause(err)
	assert.IsType(t, validator.ValidationErrors{}, cause)
	validatorErr := cause.(validator.ValidationErrors)
	assert.Len(t, validatorErr, 2)
	assert.Equal(t, "Name", validatorErr[0].Field())
	assert.Equal(t, "required", validatorErr[0].Tag())
	assert.Equal(t, "Products", validatorErr[1].Field())
	assert.Equal(t, "min", validatorErr[1].Tag())
}

func TestLoaderBindingValidate_WhenConfigIsHasEmptyField_ShouldReturnError(t *testing.T) {
	loader, err := NewLoader(Option{
		ActiveProfiles: []string{"test_validation_empty_field"},
		ConfigPaths:    []string{"./test_assets"},
		ConfigFormat:   "yml",
	}, []Properties{new(testStoreWithValidation)})
	assert.NoError(t, err)

	props := testStoreWithValidation{}
	err = loader.Bind(&props)
	assert.Error(t, err)
	cause := errors.Cause(err)
	assert.IsType(t, validator.ValidationErrors{}, cause)
	validatorErr := cause.(validator.ValidationErrors)
	assert.Len(t, validatorErr, 1)
	assert.Equal(t, "Name", validatorErr[0].Field())
	assert.Equal(t, "required", validatorErr[0].Tag())
}

func TestLoaderBindingValidate_WhenConfigContainsFieldHasDefaultButOverride_ShouldReturnError(t *testing.T) {
	loader, err := NewLoader(Option{
		ActiveProfiles: []string{"test_validation_default_field_override"},
		ConfigPaths:    []string{"./test_assets"},
		ConfigFormat:   "yml",
	}, []Properties{new(testStoreWithValidation)})
	assert.NoError(t, err)

	props := testStoreWithValidation{}
	err = loader.Bind(&props)
	assert.Error(t, err)
	cause := errors.Cause(err)
	assert.IsType(t, validator.ValidationErrors{}, cause)
	validatorErr := cause.(validator.ValidationErrors)
	assert.Len(t, validatorErr, 1)
	assert.Equal(t, "Path", validatorErr[0].Field())
	assert.Equal(t, "required", validatorErr[0].Tag())
}
