package config

import (
	"github.com/spf13/viper"
	assert "github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

type testStructs struct {
	Key1 string
	Key2 []string        `default:"[\"val1\"]"`
	Key3 *testSubStructs `default:"{}"`
	Key4 testSubStructs
	Key5 []testSubStruct2
}

func (t testStructs) Prefix() string {
	return "parent1.parent2"
}

type testSubStructs struct {
	SubKey1 string `default:"sub_val1"`
	SubKey2 int
}

type testSubStruct2 struct {
	S1 string
	S2 int
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
	err3 := os.Setenv("PARENT1_PARENT2_KEY5_0_S1", "override_s1")
	assert.NoError(t, err3)
	defer func() {
		_ = os.Unsetenv("PARENT1_PARENT2_KEY1")
		_ = os.Unsetenv("PARENT1_PARENT2_KEY3_SUBKEY2")
		_ = os.Unsetenv("PARENT1_PARENT2_KEY5_0_S1")
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
	assert.Equal(t, "override_s1", props.Key5[0].S1)
}

var yamlDeepNestedSlices = []byte(`application:
  name: WSLT Service Public
vinid:
  wslt:
    length: 32
    ttl: 2m
  security:
    http:
      jwt.type: JWT_TOKEN_MOBILE
      publicUrls:
        - /actuator/health
        - /actuator/info
        - /swagger/*
      protectedUrls:
        - { urlPattern: "/v1/tokens", method: POST, roles: [ "MOBILE_APP" ], unauthorizedWwwAuthenticateHeaders: [ "Bearer" ] }
`)

func TestSliceIndexAccess(t *testing.T) {
	v := viper.NewWithOptions()
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(keyDelimiter, "_"))
	r := strings.NewReader(string(yamlDeepNestedSlices))
	err := v.MergeConfig(r)
	assert.NoError(t, err)

	assert.Equal(t, "POST", v.GetString("vinid.security.http.protectedurls.0.method"))
}

//
//var yamlExampleWithDot = []byte(`Hacker: true
//name: steve
//kafka::topic: abc
//hobbies:
//  - skateboarding
//  - snowboarding
//  - go
//clothing:
//  jacket: leather
//  trousers: denim
//  pants:
//    size: large
//age: 35
//eyes : brown
//beard: true
//emails:
//  steve@hacker.com:
//    created: 01/02/03
//    active: true
//`)
//
//func TestKeyDelimiter(t *testing.T) {
//	v := viper.NewWithOptions(viper.KeyDelimiter("::"))
//	v.SetConfigType("yaml")
//	r := strings.NewReader(string(yamlExampleWithDot))
//	err := v.MergeConfig(r)
//	assert.NoError(t, err)
//
//	assert.Equal(t, "leather", v.GetString("kafka::topic"))
//	assert.Equal(t, "leather", v.GetString("clothing::jacket"))
//	assert.Equal(t, "01/02/03", v.GetString("emails::steve@hacker.com::created"))
//}
