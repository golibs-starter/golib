package utils

import (
	assert "github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func Test_DeepSearchInMap(t *testing.T) {
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
			name: "Give map in map When search a child map Should return correct",
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
			name: "Give map in map and child map has multiple key When search a child map Should return correct",
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
			name: "Give map in map When search with key not contains key delimiter Should return correct",
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
			name: "Give map in map When search a value Should return correct",
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
			name: "Give a simple map When search a value Should return correct",
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
			if got := DeepSearchInMap(tt.args.m, tt.args.key, "."); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeepSearchInMap() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("Give a map When run deep search Should not change input map", func(t *testing.T) {
		m := map[string]interface{}{
			"a1": map[string]interface{}{
				"a2": 1,
			},
		}
		result := DeepSearchInMap(m, "a1.a2", ".")
		assert.NotNil(t, result)
		assert.Equal(t, map[string]interface{}{
			"a1": map[string]interface{}{
				"a2": 1,
			},
		}, m)
	})
}
