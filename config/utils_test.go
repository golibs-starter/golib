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

func Test_wrapKeysAroundMap(t *testing.T) {
	type args struct {
		paths    []string
		inMap    map[string]interface{}
		endValue interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "Test when end value is not nil",
			args: args{
				paths:    []string{"k1", "k2", "k3"},
				endValue: 10,
				inMap:    nil,
			},
			want: map[string]interface{}{
				"k1": map[string]interface{}{
					"k2": map[string]interface{}{
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
			want: map[string]interface{}{
				"k1": map[string]interface{}{
					"k2": map[string]interface{}{
						"k3": nil,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := wrapKeysAroundMap(tt.args.paths, tt.args.endValue, tt.args.inMap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("wrapKeysAroundMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_deepSearchInMap(t *testing.T) {
	type args struct {
		m   map[string]interface{}
		key string
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "test 1",
			args: args{
				m: map[string]interface{}{
					"a1": map[string]interface{}{
						"a2": map[string]interface{}{
							"a3": 1,
						},
					},
				},
				key: "a1.a2",
			},
			want: map[string]interface{}{
				"a3": 1,
			},
		},
		{
			name: "test 2",
			args: args{
				m: map[string]interface{}{
					"a1": map[string]interface{}{
						"a2": map[string]interface{}{
							"a3": 1,
							"a4": 2,
						},
					},
				},
				key: "a1.a2",
			},
			want: map[string]interface{}{
				"a3": 1,
				"a4": 2,
			},
		},
		{
			name: "test 3",
			args: args{
				m: map[string]interface{}{
					"a1": map[string]interface{}{
						"a2": 1,
					},
				},
				key: "a1",
			},
			want: map[string]interface{}{
				"a2": 1,
			},
		},
		{
			name: "test 4",
			args: args{
				m: map[string]interface{}{
					"a1": map[string]interface{}{
						"a2": 1,
					},
				},
				key: "a1.a2",
			},
			want: map[string]interface{}{},
		},
		{
			name: "test 5",
			args: args{
				m: map[string]interface{}{
					"a1": 1,
				},
				key: "a1",
			},
			want: map[string]interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := deepSearchInMap(tt.args.m, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("deepSearchInMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mapToLowerKey(t *testing.T) {
	type args struct {
		mp map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "Give a map should return map with lower key",
			args: args{mp: map[string]interface{}{
				"K1": map[string]interface{}{
					"K2": 1,
				},
				"k3": map[string]interface{}{
					"K4": "v",
				},
			}},
			want: map[string]interface{}{
				"K1": map[string]interface{}{
					"K2": 1,
				},
				"k3": map[string]interface{}{
					"K4": "v",
				},
			},
		},
		{
			name: "Give a map with duplicated keys should return map with lower key and override correctly",
			args: args{mp: map[string]interface{}{
				"K1": map[string]interface{}{
					"K2": 1,
				},
				// lower case k3
				"k3": map[string]interface{}{
					"K4": "v1",
				},
				// upper case K3
				"K3": map[string]interface{}{
					"K4": "v2",
				},
			}},
			want: map[string]interface{}{
				"K1": map[string]interface{}{
					"K2": 1,
				},
				"k3": map[string]interface{}{
					"K4": "v2",
				},
			},
		},
		{
			name: "Give a nil map should return nil",
			args: args{mp: nil},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapToLowerKey(tt.args.mp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapToLowerKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
