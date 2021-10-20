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
	Address        string `mapstructure:"BuildingAddress" default:"Apple Centre Building"`
	Open           bool   `default:"true"`
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
	Images  map[string]testVariantImage
}

type testVariantImage struct {
	Size      string `default:"normal"`
	Width     int64
	Height    int64
	IsDefault bool `default:"true"`
}

func TestLoaderBinding_WhenCustomizeProps_WithInlineParent_ShouldReturnWithCorrectValue(t *testing.T) {
	err := os.Setenv("ORG_STORE_PRODUCTS_0_PRICE", "610")
	assert.NoError(t, err)
	defer func() {
		_ = os.Unsetenv("ORG_STORE_PRODUCTS_0_PRICE")
	}()

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
	assert.Equal(t, "123 Nguyen Trai, q.Thanh Xuan, tp.Ha Noi", props.Address)
	assert.Equal(t, []string{"Iphone", "Ipad"}, props.Tags)
	assert.Equal(t, []string{"0967xxx", "0968xxx"}, props.PhoneNumbers)
	assert.Equal(t, 1, props.NumberProducts)
	assert.True(t, props.Open)
	assert.Len(t, props.Products, 1)

	assert.Equal(t, "Iphone 6", props.Products[0].Title)
	assert.EqualValues(t, 610, props.Products[0].Price)
	assert.Equal(t, "$", props.Products[0].Currency)
	assert.Len(t, props.Products[0].Variants, 1)
	assert.Equal(t, "red", props.Products[0].Variants[0].Color)
	assert.Equal(t, "64gb", props.Products[0].Variants[0].Storage)
	assert.Len(t, props.Products[0].Variants[0].Images, 2)
	assert.Equal(t, "120", props.Products[0].Variants[0].Images["Sm"].Size)
	assert.Equal(t, int64(120), props.Products[0].Variants[0].Images["Sm"].Width)
	assert.Equal(t, int64(80), props.Products[0].Variants[0].Images["Sm"].Height)
	assert.Equal(t, "800", props.Products[0].Variants[0].Images["Xl"].Size)
	assert.Equal(t, int64(800), props.Products[0].Variants[0].Images["Xl"].Width)
	assert.Equal(t, int64(600), props.Products[0].Variants[0].Images["Xl"].Height)
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
	assert.Len(t, props.Products, 1)

	assert.Equal(t, "Iphone 6", props.Products[0].Title)
	assert.EqualValues(t, 600, props.Products[0].Price)
	assert.Equal(t, "$", props.Products[0].Currency)
	assert.Len(t, props.Products[0].Variants, 1)
	assert.Equal(t, "red", props.Products[0].Variants[0].Color)
	assert.Equal(t, "64gb", props.Products[0].Variants[0].Storage)
	assert.Len(t, props.Products[0].Variants[0].Images, 2)
	assert.Equal(t, "120", props.Products[0].Variants[0].Images["Sm"].Size)
	assert.Equal(t, int64(120), props.Products[0].Variants[0].Images["Sm"].Width)
	assert.Equal(t, int64(80), props.Products[0].Variants[0].Images["Sm"].Height)
	assert.Equal(t, "800", props.Products[0].Variants[0].Images["Xl"].Size)
	assert.Equal(t, int64(800), props.Products[0].Variants[0].Images["Xl"].Width)
	assert.Equal(t, int64(600), props.Products[0].Variants[0].Images["Xl"].Height)
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

func TestLoaderBinding_WhenDefaultHasBeenSet_ShouldReturnWithCorrectDefaultValueInSlice(t *testing.T) {
	err := os.Setenv("ORG_STORE_PRODUCTS_0_PRICE", "610")
	assert.NoError(t, err)
	defer func() {
		_ = os.Unsetenv("ORG_STORE_PRODUCTS_0_PRICE")
	}()

	loader, err := NewLoader(Option{
		ActiveProfiles: []string{"test_default_config"},
		ConfigPaths:    []string{"./test_assets"},
		ConfigFormat:   "yaml",
	}, []Properties{new(testStore)})
	assert.NoError(t, err)

	props := testStore{}
	err = loader.Bind(&props)
	assert.NoError(t, err)

	assert.Equal(t, "Apple", props.Name)
	assert.Equal(t, "Hanoi", props.Location)
	assert.Equal(t, "Apple Centre Building", props.Address)
	assert.Equal(t, []string{"Iphone", "Ipad"}, props.Tags)
	assert.Equal(t, []string{"0967xxx", "0968xxx"}, props.PhoneNumbers)
	assert.False(t, props.Open)
	assert.Equal(t, 1, props.NumberProducts)
	assert.Len(t, props.Products, 1)

	assert.Equal(t, "Iphone 6", props.Products[0].Title)
	assert.EqualValues(t, 610, props.Products[0].Price)
	assert.Equal(t, "$", props.Products[0].Currency)
	assert.Len(t, props.Products[0].Variants, 1)
	assert.Equal(t, "red", props.Products[0].Variants[0].Color)
	assert.Equal(t, "64gb", props.Products[0].Variants[0].Storage)
	assert.Len(t, props.Products[0].Variants[0].Images, 2)
	assert.EqualValues(t, 120, props.Products[0].Variants[0].Images["Normal"].Width)
	assert.EqualValues(t, 80, props.Products[0].Variants[0].Images["Normal"].Height)
	assert.Equal(t, "", props.Products[0].Variants[0].Images["Normal"].Size)
	assert.True(t, props.Products[0].Variants[0].Images["Normal"].IsDefault)
	assert.EqualValues(t, 80, props.Products[0].Variants[0].Images["Thumb"].Width)
	assert.EqualValues(t, 80, props.Products[0].Variants[0].Images["Thumb"].Height)
	assert.Equal(t, "thumb", props.Products[0].Variants[0].Images["Thumb"].Size)
	assert.False(t, props.Products[0].Variants[0].Images["Thumb"].IsDefault)
}

func TestLoaderBinding_WhenConfigWithPlaceholderValue_AndEnvHasBeenSet_ShouldReturnWithValueInEnv(t *testing.T) {
	err1 := os.Setenv("STORE_LOCATION", "Haiduong")
	assert.NoError(t, err1)
	err2 := os.Setenv("PRICE_CURRENCY", "Dolar")
	assert.NoError(t, err2)
	err3 := os.Setenv("PREMIUM_BLUE_STORAGE", "1Tb")
	assert.NoError(t, err3)
	defer func() {
		_ = os.Unsetenv("STORE_LOCATION")
		_ = os.Unsetenv("PRICE_CURRENCY")
		_ = os.Unsetenv("PREMIUM_BLUE_STORAGE")
	}()

	loader, err := NewLoader(Option{
		ActiveProfiles: []string{"test_placeholder_values"},
		ConfigPaths:    []string{"./test_assets"},
		ConfigFormat:   "yaml",
	}, []Properties{new(testStore)})
	assert.NoError(t, err)

	props := testStore{}
	err = loader.Bind(&props)
	assert.NoError(t, err)

	assert.Equal(t, "Apple", props.Name)
	assert.Equal(t, "Haiduong", props.Location)
	assert.Len(t, props.Products, 1)

	assert.Equal(t, "Iphone 6", props.Products[0].Title)
	assert.EqualValues(t, 600, props.Products[0].Price)
	assert.Equal(t, "Dolar", props.Products[0].Currency)

	assert.Len(t, props.Products[0].Variants, 2)
	assert.Equal(t, "red", props.Products[0].Variants[0].Color)
	assert.Equal(t, "64gb", props.Products[0].Variants[0].Storage)
	assert.Equal(t, "premium_blue", props.Products[0].Variants[1].Color)
	assert.Equal(t, "1Tb", props.Products[0].Variants[1].Storage)
}

func TestLoaderBinding_WhenConfigWithPlaceholderValue_AndEnvIsNotSet_ShouldReturnError(t *testing.T) {
	err1 := os.Setenv("STORE_LOCATION", "Haiduong")
	assert.NoError(t, err1)
	err2 := os.Setenv("PRICE_CURRENCY", "Dolar")
	assert.NoError(t, err2)
	// Missing PREMIUM_BLUE_STORAGE
	defer func() {
		_ = os.Unsetenv("STORE_LOCATION")
		_ = os.Unsetenv("PRICE_CURRENCY")
	}()

	loader, err := NewLoader(Option{
		ActiveProfiles: []string{"test_placeholder_values"},
		ConfigPaths:    []string{"./test_assets"},
		ConfigFormat:   "yaml",
	}, []Properties{new(testStore)})
	assert.NoError(t, err)

	props := testStore{}
	err = loader.Bind(&props)
	assert.Error(t, err)
}

func TestLoaderBinding_WhenOverrideByKeyWithCaseInsensitive_ShouldReturnCorrect(t *testing.T) {
	loader, err := NewLoader(Option{
		ActiveProfiles: []string{"test_key_case_insensitive"},
		ConfigPaths:    []string{"./test_assets"},
		ConfigFormat:   "yaml",
	}, []Properties{new(testStore)})
	assert.NoError(t, err)

	props := testStore{}
	err = loader.Bind(&props)
	assert.NoError(t, err)
	assert.Equal(t, "Apple Inc", props.Name)
}

func TestLoaderBinding_WhenOverrideByMultipleKeysWithCaseInsensitive_ShouldReturnCorrect(t *testing.T) {
	loader, err := NewLoader(Option{
		ActiveProfiles: []string{"test_key_case_insensitive_override"},
		ConfigPaths:    []string{"./test_assets"},
		ConfigFormat:   "yaml",
	}, []Properties{new(testStore)})
	assert.NoError(t, err)

	props := testStore{}
	err = loader.Bind(&props)
	assert.NoError(t, err)
	assert.Equal(t, "Apple Company", props.Name)
	assert.Equal(t, "Hanoi City", props.Location)
}

func TestLoaderBinding_WhenProfileFileInYamlFormat_ShouldReturnCorrect(t *testing.T) {
	loader, err := NewLoader(Option{
		ActiveProfiles: []string{"file_in_yaml_format"},
		ConfigPaths:    []string{"./test_assets"},
		ConfigFormat:   "yml",
	}, []Properties{new(testStore)})
	assert.NoError(t, err)

	props := testStore{}
	err = loader.Bind(&props)
	assert.NoError(t, err)
	assert.Equal(t, "Apple Inc.", props.Name)
}
