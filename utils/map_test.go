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

func Test_mergeMaps(t *testing.T) {
	type args struct {
		src map[string]interface{}
		tgt map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "when modify map item and slice item should merge correct",
			args: args{
				src: map[string]interface{}{
					"k1": "v1-modified",
					"k2": map[string]interface{}{
						"k22": "v22-modified",
					},
					"k3": []interface{}{
						map[string]interface{}{
							"k31": "v31-modified",
						},
						map[string]interface{}{ // test slice item with different type
							"k33": "v33",
						},
						map[string]interface{}{ // append new value to slice
							"k34": "v34",
							"k35": 35,
						},
						"v36-with-diff-type", // test append item with different type
					},
					"k4": 4, // test override when different type
				},
				tgt: map[string]interface{}{
					"k1": "v1",
					"k2": map[string]interface{}{
						"k21": "v21",
						"k22": "v22",
					},
					"k3": []interface{}{
						map[string]interface{}{
							"k31": "v31",
							"k32": "v32",
						},
						"v33-diff-type-item",
					},
					"k4": "v4",
				},
			},
			want: map[string]interface{}{
				"k1": "v1-modified",
				"k2": map[string]interface{}{
					"k21": "v21",
					"k22": "v22-modified",
				},
				"k3": []interface{}{
					map[string]interface{}{
						"k31": "v31-modified",
						"k32": "v32",
					},
					map[string]interface{}{
						"k33": "v33",
					},
					map[string]interface{}{
						"k34": "v34",
						"k35": 35,
					},
					"v36-with-diff-type",
				},
				"k4": 4,
			},
		},
		{
			name: "when map item is map[interface{}]interface{} should merge and convert to map[string]interface{}",
			args: args{
				src: map[string]interface{}{
					"k1": "v1-modified",
					"k2": map[interface{}]interface{}{
						"k22": "v22-modified",
					},
					"k3": []interface{}{
						map[string]interface{}{
							"k31": "v31-modified",
						},
						map[interface{}]interface{}{ // test slice item with different type
							"k33": "v33",
						},
						map[interface{}]interface{}{ // append new value to slice
							"k34": "v34",
							"k35": 35,
						},
					},
				},
				tgt: map[string]interface{}{
					"k1": "v1",
					"k2": map[string]interface{}{
						"k21": "v21",
						"k22": "v22",
					},
					"k3": []interface{}{
						map[string]interface{}{
							"k31": "v31",
							"k32": "v32",
						},
						"v33-diff-type-item",
					},
				},
			},
			want: map[string]interface{}{
				"k1": "v1-modified",
				"k2": map[string]interface{}{
					"k21": "v21",
					"k22": "v22-modified",
				},
				"k3": []interface{}{
					map[string]interface{}{
						"k31": "v31-modified",
						"k32": "v32",
					},
					map[string]interface{}{
						"k33": "v33",
					},
					map[string]interface{}{
						"k34": "v34",
						"k35": 35,
					},
				},
			},
		},
		{
			name: "when map contains lower and uppercase should merge with case insensitive",
			args: args{
				src: map[string]interface{}{
					"K1": "v1-modified", // key uppercase
					"k2": map[interface{}]interface{}{
						"k22": "v22-modified",
					},
					"k3": []interface{}{
						map[string]interface{}{
							"K31": "v31-modified", // key uppercase
						},
					},
				},
				tgt: map[string]interface{}{
					"k1": "v1",
					"k2": map[string]interface{}{
						"k21": "v21",
						"K22": "v22", // key uppercase
					},
					"k3": []interface{}{
						map[string]interface{}{
							"k31": "v31",
							"k32": "v32",
						},
					},
				},
			},
			want: map[string]interface{}{
				"k1": "v1-modified",
				"k2": map[string]interface{}{
					"k21": "v21",
					"K22": "v22-modified",
				},
				"k3": []interface{}{
					map[string]interface{}{
						"k31": "v31-modified",
						"k32": "v32",
					},
				},
			},
		},
		{
			name: "when modify map inside slice should merge correct",
			args: args{
				src: map[string]interface{}{
					"k1": []interface{}{
						map[string]interface{}{
							"k11": map[interface{}]interface{}{
								"k111": "v111-modified",
							},
						},
					},
				},
				tgt: map[string]interface{}{
					"k1": []interface{}{
						map[string]interface{}{
							"k11": map[interface{}]interface{}{
								"k111": "v111",
							},
						},
					},
				},
			},
			want: map[string]interface{}{
				"k1": []interface{}{
					map[string]interface{}{
						"k11": map[string]interface{}{
							"k111": "v111-modified",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MergeCaseInsensitiveMaps(tt.args.src, tt.args.tgt)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeCaseInsensitiveMaps() = %v, want %v", tt.args.tgt, tt.want)
			}
		})
	}
}
