package config

import (
	assert "github.com/stretchr/testify/require"
	"os"
	"testing"
)

type testStructs struct {
	Key1 string
	Key2 []string        `default:"[\"val1\"]"`
	Key3 *testSubStructs `default:"{}"`
	Key4 testSubStructs
}

func (t testStructs) Prefix() string {
	return "parent1.parent2"
}

type testSubStructs struct {
	SubKey1 string `default:"sub_val1"`
	SubKey2 int
}

func TestLoaderBind_WhenNoCustomizedProps_ShouldReturnWithDefaultValue(t *testing.T) {
	loader, err := NewLoader(Option{
		ConfigPaths:  []string{"./test_assets"},
		ConfigFormat: "yaml",
	}, []Properties{new(testStructs)})
	assert.NoError(t, err)
	props := testStructs{}
	err = loader.Bind(&props)
	assert.NoError(t, err)
	assert.Equal(t, "", props.Key1)
	assert.Equal(t, []string{"val1"}, props.Key2)
	assert.NotNil(t, props.Key3)
	assert.Equal(t, "sub_val1", props.Key3.SubKey1)
	assert.Equal(t, 0, props.Key3.SubKey2)
	assert.Equal(t, "sub_val1", props.Key4.SubKey1)
	assert.Equal(t, 0, props.Key4.SubKey2)
}

func TestLoaderBind_WhenCustomizeProps_WithInlineParent_ShouldReturnWithCorrectValue(t *testing.T) {
	loader, err := NewLoader(Option{
		ActiveProfiles: []string{"test_inline_key"},
		ConfigPaths:    []string{"./test_assets"},
		ConfigFormat:   "yaml",
	}, []Properties{new(testStructs)})
	assert.NoError(t, err)

	props := testStructs{}
	err = loader.Bind(&props)
	assert.NoError(t, err)

	assert.Equal(t, "key1_val", props.Key1)
	assert.Equal(t, []string{"key2_val1", "key2_val2"}, props.Key2)
	assert.NotNil(t, props.Key3)
	assert.Equal(t, "new_sub_val1", props.Key3.SubKey1)
	assert.Equal(t, 0, props.Key3.SubKey2)
	assert.Equal(t, "", props.Key4.SubKey1)
	assert.Equal(t, 3, props.Key4.SubKey2)
}

func TestLoaderBind_WhenCustomizeProps_WithNestedParent_ShouldReturnWithCorrectValue(t *testing.T) {
	loader, err := NewLoader(Option{
		ActiveProfiles: []string{"test_nested_key"},
		ConfigPaths:    []string{"./test_assets"},
		ConfigFormat:   "yaml",
	}, []Properties{new(testStructs)})
	assert.NoError(t, err)

	props := testStructs{}
	err = loader.Bind(&props)
	assert.NoError(t, err)

	assert.Equal(t, "key1_val", props.Key1)
	assert.Equal(t, []string{"key2_val1", "key2_val2"}, props.Key2)
	assert.NotNil(t, props.Key3)
	assert.Equal(t, "new_sub_val1", props.Key3.SubKey1)
	assert.Equal(t, 0, props.Key3.SubKey2)
	assert.Equal(t, "", props.Key4.SubKey1)
	assert.Equal(t, 3, props.Key4.SubKey2)
}

func TestLoaderBind_WhenCustomizeProps_WithNestedParentAndOverrideByAnotherFileConfig_ShouldReturnWithCorrectValue(t *testing.T) {
	loader, err := NewLoader(Option{
		ActiveProfiles: []string{"test_nested_key", "test_nested_key_1"},
		ConfigPaths:    []string{"./test_assets"},
		ConfigFormat:   "yaml",
	}, []Properties{new(testStructs)})
	assert.NoError(t, err)

	props := testStructs{}
	err = loader.Bind(&props)
	assert.NoError(t, err)

	assert.Equal(t, "overwrite_val1_by_file", props.Key1)
	assert.Equal(t, []string{"key2_val1", "key2_val2"}, props.Key2)
	assert.NotNil(t, props.Key3)
	assert.Equal(t, "new_sub_val1", props.Key3.SubKey1)
	assert.Equal(t, 0, props.Key3.SubKey2)
	assert.Equal(t, "", props.Key4.SubKey1)
	assert.Equal(t, 3, props.Key4.SubKey2)
}

func TestLoaderBind_WhenCustomizeProps_WithNestedParentAndEnvHasBeenSet_ShouldReturnWithCorrectValue(t *testing.T) {
	err1 := os.Setenv("PARENT1_PARENT2_KEY1", "test_override_val1")
	assert.NoError(t, err1)
	err2 := os.Setenv("PARENT1_PARENT2_KEY3_SUBKEY2", "18")
	assert.NoError(t, err2)
	defer func() {
		_ = os.Unsetenv("PARENT1_PARENT2_KEY1")
		_ = os.Unsetenv("PARENT1_PARENT2_KEY3_SUBKEY2")
	}()

	loader, err := NewLoader(Option{
		ActiveProfiles: []string{"test_nested_key"},
		ConfigPaths:    []string{"./test_assets"},
		ConfigFormat:   "yaml",
	}, []Properties{new(testStructs)})
	assert.NoError(t, err)

	props := testStructs{}
	err = loader.Bind(&props)
	assert.NoError(t, err)

	assert.Equal(t, "test_override_val1", props.Key1)
	assert.Equal(t, []string{"key2_val1", "key2_val2"}, props.Key2)
	assert.NotNil(t, props.Key3)
	assert.Equal(t, "new_sub_val1", props.Key3.SubKey1)
	assert.Equal(t, 18, props.Key3.SubKey2)
	assert.Equal(t, "", props.Key4.SubKey1)
	assert.Equal(t, 3, props.Key4.SubKey2)
}
