package config

import (
	assert "github.com/stretchr/testify/require"
	"os"
	"testing"
)

type testStore struct {
	Name           string
	Location       string
	Tags           []string
	PhoneNumbers   []string `default:"[\"0967xxx\", \"0968xxx\"]"`
	NumberProducts int
	Products       []testProduct
}

func (t testStore) Prefix() string {
	return "org.store"
}

type testProduct struct {
	Title    string
	Price    int64
	Currency string `default:"$"`
	Variants []*testVariant
}

type testVariant struct {
	Color   string
	Storage string
}

func TestLoaderBinding_WhenCustomizeProps_WithInlineParent_ShouldReturnWithCorrectValue(t *testing.T) {
	err := os.Setenv("ORG_STORE_PRODUCTS_0_PRICE", "610")
	assert.NoError(t, err)

	loader, err := NewLoader(Option{
		ActiveProfiles: []string{"test_inline_key"},
		ConfigPaths:    []string{"./test_assets"},
		ConfigFormat:   "yaml",
	}, []Properties{new(testStore)})
	assert.NoError(t, err)

	props := testStore{}
	err = loader.Bind(&props)
	assert.NoError(t, err)

	assert.Equal(t, "Apple", props.Name)
	assert.Equal(t, "Hanoi", props.Location)
	assert.Equal(t, []string{"Iphone", "Ipad"}, props.Tags)
	assert.Equal(t, []string{"0967xxx", "0968xxx"}, props.PhoneNumbers)
	assert.Equal(t, 1, props.NumberProducts)
	assert.Len(t, props.Products, 1)

	assert.Equal(t, "Iphone 6", props.Products[0].Title)
	assert.EqualValues(t, 610, props.Products[0].Price)
	assert.Equal(t, "$", props.Products[0].Currency)
	assert.Len(t, props.Products[0].Variants, 1)
	assert.Equal(t, "red", props.Products[0].Variants[0].Color)
	assert.Equal(t, "64gb", props.Products[0].Variants[0].Storage)
}

func TestLoaderBinding_WhenCustomizeProps_WithInlineKeyOverrideNestedKey_ShouldReturnWithCorrectValue(t *testing.T) {
	loader, err := NewLoader(Option{
		ActiveProfiles: []string{"test_inline_key_override"}, //override value in default profile
		ConfigPaths:    []string{"./test_assets"},
		ConfigFormat:   "yaml",
	}, []Properties{new(testStore)})
	assert.NoError(t, err)

	props := testStore{}
	err = loader.Bind(&props)
	assert.NoError(t, err)

	assert.Equal(t, "Apple Store", props.Name)
	assert.Equal(t, "Hanoi", props.Location)
	assert.Len(t, props.Tags, 0)
	assert.Equal(t, []string{"0967xxx", "0968xxx"}, props.PhoneNumbers)
}

// TODO fix code to cover this test, currently inline key is cannot override by nested key
func TestLoaderBinding_WhenCustomizeProps_WithNestedKeyOverrideInlineKey_ShouldReturnWithCorrectValue(t *testing.T) {
	loader, err := NewLoader(Option{
		ActiveProfiles: []string{"test_inline_key", "test_nested_key_override"},
		ConfigPaths:    []string{"./test_assets"},
		ConfigFormat:   "yaml",
	}, []Properties{new(testStore)})
	assert.NoError(t, err)

	props := testStore{}
	err = loader.Bind(&props)
	assert.NoError(t, err)

	assert.Equal(t, "Apple", props.Name)
	assert.Equal(t, "Vietnam", props.Location)
	assert.Equal(t, []string{"iphone"}, props.Tags)
	assert.Equal(t, []string{"0969xxx", "0970xxx"}, props.PhoneNumbers)
	assert.Equal(t, 1, props.NumberProducts)
}

func TestLoaderBinding_WhenCustomizeProps_AndEnvHasBeenSet_ShouldReturnWithCorrectValue(t *testing.T) {
	err1 := os.Setenv("ORG_STORE_NUMBERPRODUCTS", "3")
	assert.NoError(t, err1)
	err2 := os.Setenv("ORG_STORE_PRODUCTS_0_PRICE", "610")
	assert.NoError(t, err2)
	err3 := os.Setenv("ORG_STORE_PRODUCTS_0_VARIANTS_1_COLOR", "space_blue")
	assert.NoError(t, err3)
	err4 := os.Setenv("ORG_STORE_PRODUCTS_1_TITLE", "Iphone 13 Pro Max")
	assert.NoError(t, err4)
	defer func() {
		_ = os.Unsetenv("ORG_STORE_NUMBERPRODUCTS")
		_ = os.Unsetenv("ORG_STORE_PRODUCTS_0_PRICE")
		_ = os.Unsetenv("ORG_STORE_PRODUCTS_0_VARIANTS_1_COLOR")
		_ = os.Unsetenv("ORG_STORE_PRODUCTS_1_TITLE")
	}()

	loader, err := NewLoader(Option{
		ActiveProfiles: []string{"test_nested_key"},
		ConfigPaths:    []string{"./test_assets"},
		ConfigFormat:   "yaml",
	}, []Properties{new(testStore)})
	assert.NoError(t, err)

	props := testStore{}
	err = loader.Bind(&props)
	assert.NoError(t, err)

	assert.Equal(t, "Apple", props.Name)
	assert.Equal(t, "Hanoi", props.Location)
	assert.Equal(t, []string{"Iphone", "Ipad"}, props.Tags)
	assert.Equal(t, []string{"0967xxx", "0968xxx"}, props.PhoneNumbers)
	assert.Equal(t, 3, props.NumberProducts)
	assert.Len(t, props.Products, 2)

	assert.Equal(t, "Iphone 6", props.Products[0].Title)
	assert.EqualValues(t, 610, props.Products[0].Price)
	assert.Equal(t, "$", props.Products[0].Currency)
	assert.Len(t, props.Products[0].Variants, 2)
	assert.Equal(t, "red", props.Products[0].Variants[0].Color)
	assert.Equal(t, "64gb", props.Products[0].Variants[0].Storage)
	assert.Equal(t, "space_blue", props.Products[0].Variants[1].Color)
	assert.Equal(t, "128gb", props.Products[0].Variants[1].Storage)

	assert.Equal(t, "Iphone 13 Pro Max", props.Products[1].Title)
	assert.EqualValues(t, 45000000, props.Products[1].Price)
	assert.Equal(t, "VND", props.Products[1].Currency)
	assert.Len(t, props.Products[1].Variants, 1)
	assert.Equal(t, "yellow", props.Products[1].Variants[0].Color)
	assert.Equal(t, "1TB", props.Products[1].Variants[0].Storage)
}

// TODO fix code to cover this test, currently default value inside slice is not working
func TestLoaderBinding_WhenNoCustomizeValueInSlice_ShouldReturnWithCorrectDefaultValueInSlice(t *testing.T) {
	err := os.Setenv("ORG_STORE_PRODUCTS_0_PRICE", "610")
	assert.NoError(t, err)

	loader, err := NewLoader(Option{
		ActiveProfiles: []string{"test_default_in_slice"},
		ConfigPaths:    []string{"./test_assets"},
		ConfigFormat:   "yaml",
	}, []Properties{new(testStore)})
	assert.NoError(t, err)

	props := testStore{}
	err = loader.Bind(&props)
	assert.NoError(t, err)

	assert.Equal(t, "Apple", props.Name)
	assert.Equal(t, "Hanoi", props.Location)
	assert.Equal(t, []string{"Iphone", "Ipad"}, props.Tags)
	assert.Equal(t, []string{"0967xxx", "0968xxx"}, props.PhoneNumbers)
	assert.Equal(t, 1, props.NumberProducts)
	assert.Len(t, props.Products, 1)

	assert.Equal(t, "Iphone 6", props.Products[0].Title)
	assert.EqualValues(t, 610, props.Products[0].Price)
	assert.Equal(t, "$", props.Products[0].Currency)
	assert.Len(t, props.Products[0].Variants, 1)
	assert.Equal(t, "red", props.Products[0].Variants[0].Color)
	assert.Equal(t, "64gb", props.Products[0].Variants[0].Storage)
}

// TODO add test to cover placeholder value
