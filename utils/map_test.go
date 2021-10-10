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

func Test_WrapKeysAroundMap(t *testing.T) {
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
			name: "Give a paths When end value is not nil Should return correct",
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
			name: "Give a paths When end value is nil Should return correct",
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
			if got := WrapKeysAroundMap(tt.args.paths, tt.args.endValue, tt.args.inMap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WrapKeysAroundMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
