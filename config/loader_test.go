package config

import (
	assert "github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestLoader_WhenFormatIsNotSupported_ShouldReturnError(t *testing.T) {
	_, err := NewLoader(Option{
		ActiveProfiles: []string{"file_in_yaml_format"},
		ConfigPaths:    []string{"./test_assets"},
		ConfigFormat:   "ini",
	}, []Properties{new(testStore)})
	assert.ErrorIs(t, err, ErrFormatNotSupported)
}

func Test_buildEnvKeys(t *testing.T) {
	type args struct {
		data                                 map[interface{}]interface{}
		keyDelim, envDelim, baseKey, baseEnv string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "given mixed data type should return correct",
			args: args{
				data: map[interface{}]interface{}{
					"k1": "v1",
					"k2": map[string]interface{}{
						"k21": "v21",
						"k22": map[interface{}]interface{}{},
					},
					"k3": map[interface{}]interface{}{
						"k31": map[string]interface{}{
							"k311": "v311",
						},
					},
				},
				keyDelim: ".",
				envDelim: "_",
				baseKey:  "",
				baseEnv:  "",
			},
			want: map[string]string{
				"k1":          "K1",
				"k2.k21":      "K2_K21",
				"k2.k22":      "K2_K22",
				"k3.k31.k311": "K3_K31_K311",
			},
		},
		{
			name: "given base key and env should return correct",
			args: args{
				data: map[interface{}]interface{}{
					"k1": "v1",
					"k2": map[string]interface{}{
						"k21": "v21",
					},
				},
				keyDelim: ",",
				envDelim: "-",
				baseKey:  "base1",
				baseEnv:  "base2",
			},
			want: map[string]string{
				"base1,k1":     "BASE2-K1",
				"base1,k2,k21": "BASE2-K2-K21",
			},
		},
		{
			name: "given data is nil should return empty",
			args: args{
				data: nil,
			},
			want: map[string]string{},
		},
		{
			name: "given data is empty should return empty",
			args: args{
				data: map[interface{}]interface{}{},
			},
			want: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildEnvKeys(tt.args.data, tt.args.keyDelim, tt.args.envDelim, tt.args.baseKey, tt.args.baseEnv); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildEnvKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}
