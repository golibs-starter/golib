package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"testing"
)

func Test_ReplacePlaceholderValue_WhenValidPlaceholderAndEnvIsPresent_ShouldReturnCorrect(t *testing.T) {
	_ = os.Setenv("ENV_EXAMPLE", "test")
	defer func() {
		_ = os.Unsetenv("ENV_EXAMPLE")
	}()
	val, err := ReplacePlaceholderValue("${ENV_EXAMPLE}")
	assert.Nil(t, err)
	assert.Equal(t, "test", val)
}

func Test_ReplacePlaceholderValue_WhenValueIsNotString_ShouldReturnCurrentValue(t *testing.T) {
	val, err := ReplacePlaceholderValue(10)
	assert.Nil(t, err)
	assert.Equal(t, 10, val)
}

func Test_ReplacePlaceholderValue_WhenItIsNotPlaceholder_ShouldReturnCurrentValue(t *testing.T) {
	val1, err := ReplacePlaceholderValue("TEST}")
	assert.Nil(t, err)
	assert.Equal(t, "TEST}", val1)

	val2, err := ReplacePlaceholderValue("${TEST")
	assert.Nil(t, err)
	assert.Equal(t, "${TEST", val2)

	val3, err := ReplacePlaceholderValue(" ${TEST}") //starts with space
	assert.Nil(t, err)
	assert.Equal(t, " ${TEST}", val3)
}

func Test_ReplacePlaceholderValue_WhenEmptyPlaceholderKey_ShouldReturnError(t *testing.T) {
	val, err := ReplacePlaceholderValue("${}")
	assert.NotNil(t, err)
	assert.Nil(t, val)
}

func Test_ReplacePlaceholderValue_WhenValidPlaceholderAndEnvNotPresent_ShouldReturnError(t *testing.T) {
	val, err := ReplacePlaceholderValue("${ENV_EXAMPLE}")
	assert.NotNil(t, err)
	assert.Nil(t, val)
}

func Test_ReplacePlaceholderValue_WhenValidPlaceholderAndEnvIsPresentAndEmpty_ShouldReturnEmptyValue(t *testing.T) {
	_ = os.Setenv("ENV_EXAMPLE", "")
	defer func() {
		_ = os.Unsetenv("ENV_EXAMPLE")
	}()
	val, err := ReplacePlaceholderValue("${ENV_EXAMPLE}")
	assert.Nil(t, err)
	assert.Equal(t, "", val)
}

func Test_convertInlineKeyToMap(t *testing.T) {
	type args struct {
		paths    []string
		inMap    map[interface{}]interface{}
		endValue interface{}
	}
	tests := []struct {
		name string
		args args
		want map[interface{}]interface{}
	}{
		{
			name: "Test when end value is not nil",
			args: args{
				paths:    []string{"k1", "k2", "k3"},
				endValue: 10,
				inMap:    nil,
			},
			want: map[interface{}]interface{}{
				"k1": map[interface{}]interface{}{
					"k2": map[interface{}]interface{}{
						"k3": 10,
					},
				},
			},
		},
		{
			name: "Test when end value is nil",
			args: args{
				paths:    []string{"k1", "k2", "k3"},
				endValue: nil,
				inMap:    nil,
			},
			want: map[interface{}]interface{}{
				"k1": map[interface{}]interface{}{
					"k2": map[interface{}]interface{}{
						"k3": nil,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertSliceToNestedMap(tt.args.paths, tt.args.endValue, tt.args.inMap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertSliceToNestedMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
