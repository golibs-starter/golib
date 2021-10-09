package utils

import (
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"reflect"
	"testing"
	"time"
)

func TestLinkedHMapToMapStr(t *testing.T) {
	type args struct {
		hMap *linkedhashmap.Map
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "Give a empty LinkedHashMap Should return empty map",
			args: args{
				hMap: LinkedHMap(),
			},
			want: map[string]interface{}{},
		},
		{
			name: "Give a simple LinkedHashMap Should return correct",
			args: args{
				hMap: LinkedHMap(
					NewMapItem("a", 1),
					NewMapItem("b", 2),
				),
			},
			want: map[string]interface{}{
				"a": 1,
				"b": 2,
			},
		},
		{
			name: "Give a complex LinkedHashMap Should return correct",
			args: args{
				hMap: LinkedHMap(
					NewMapItem("a", 1),
					NewMapItem("b", LinkedHMap(
						NewMapItem("c", 2),
						NewMapItem("d", LinkedHMap(
							NewMapItem("e", 3),
							NewMapItem("f", "vf"),
							NewMapItem("a", 3*time.Minute),
						)),
					)),
					NewMapItem("h", []*linkedhashmap.Map{
						LinkedHMap(NewMapItem("i", 5)),
						LinkedHMap(NewMapItem("k", LinkedHMap(NewMapItem("n", 7)))),
					}),
				),
			},
			want: map[string]interface{}{
				"a": 1,
				"b": map[string]interface{}{
					"c": 2,
					"d": map[string]interface{}{
						"e": 3,
						"f": "vf",
						"a": 3 * time.Minute,
					},
				},
				"h": []interface{}{
					map[string]interface{}{
						"i": 5,
					},
					map[string]interface{}{
						"k": map[string]interface{}{
							"n": 7,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LinkedHMapToMapStr(tt.args.hMap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LinkedHMapToMapStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
