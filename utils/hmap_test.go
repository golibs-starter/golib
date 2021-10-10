package utils

import (
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	"reflect"
	"testing"
	"time"
)

func Test_LinkedHMapToMapStr(t *testing.T) {
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

func Test_MergeLinkedHashMap(t *testing.T) {
	t.Run("Give Dst & Src Should Return Success", func(t *testing.T) {
		dst := linkedhashmap.New()
		dst.Put("a", 1)
		dst.Put("b", "v2")

		c := linkedhashmap.New()
		c.Put("d", 3)
		c.Put("e", 4)
		dst.Put("c", c)

		src := linkedhashmap.New()
		src.Put("b", "v2-modified")
		cModified := linkedhashmap.New()
		cModified.Put("e", 5)
		cModified.Put("f", "v6")
		src.Put("c", cModified)
		src.Put("g", 7)

		MergeLinkedHMap(dst, src)
		require.Equal(t, 4, dst.Size())

		v, found := dst.Get("a")
		require.True(t, found)
		require.EqualValues(t, 1, v)

		v, found = dst.Get("b")
		require.True(t, found)
		require.EqualValues(t, "v2-modified", v)

		v, found = dst.Get("c")
		require.True(t, found)
		require.IsType(t, &linkedhashmap.Map{}, v)

		vc, found := v.(*linkedhashmap.Map).Get("d")
		require.True(t, found)
		require.EqualValues(t, 3, vc)

		vc, found = v.(*linkedhashmap.Map).Get("e")
		require.True(t, found)
		require.EqualValues(t, 5, vc)

		vc, found = v.(*linkedhashmap.Map).Get("f")
		require.True(t, found)
		require.EqualValues(t, "v6", vc)

		v, found = dst.Get("g")
		require.True(t, found)
		require.EqualValues(t, 7, v)
	})
}

func Test_ExpandInlineKeyInMap(t *testing.T) {
	assertExist := func(t *testing.T, m *linkedhashmap.Map, key string) interface{} {
		v, f := m.Get(key)
		require.True(t, f)
		return v
	}

	assertChildMap := func(t *testing.T, m *linkedhashmap.Map, key string) *linkedhashmap.Map {
		v := assertExist(t, m, key)
		require.IsType(t, &linkedhashmap.Map{}, v)
		return v.(*linkedhashmap.Map)
	}

	assertEquals := func(t *testing.T, m *linkedhashmap.Map, key string, val interface{}) {
		v := assertExist(t, m, key)
		require.Equal(t, val, v)
	}

	t.Run("Give inline key Should expand success", func(t *testing.T) {
		got := ExpandInlineKeyInLinkedHMap(LinkedHMap(
			NewMapItem("a.b", LinkedHMap(NewMapItem("c", 1))),
			NewMapItem("x", LinkedHMap(NewMapItem("y", []int{2, 3}))),
		), ".")
		require.Equal(t, 2, got.Size())
		v := assertChildMap(t, got, "a")
		v1 := assertChildMap(t, v, "b")
		v2 := assertExist(t, v1, "c")
		require.EqualValues(t, 1, v2)

		v = assertChildMap(t, got, "x")
		v3 := assertExist(t, v, "y")
		require.Equal(t, []int{2, 3}, v3)
	})

	t.Run("Give multiple inline keys Should expand success", func(t *testing.T) {
		got := ExpandInlineKeyInLinkedHMap(LinkedHMap(
			NewMapItem("a.b.c", 1),
			NewMapItem("x.y", []int{2, 3}),
		), ".")

		require.Equal(t, 2, got.Size())
		v := assertChildMap(t, got, "a")
		v1 := assertChildMap(t, v, "b")
		v2 := assertExist(t, v1, "c")
		require.EqualValues(t, 1, v2)

		v = assertChildMap(t, got, "x")
		v2 = assertExist(t, v, "y")
		require.Equal(t, []int{2, 3}, v2)
	})

	t.Run("Give a override key Should return success", func(t *testing.T) {
		actual := ExpandInlineKeyInLinkedHMap(LinkedHMap(
			NewMapItem("a.b", LinkedHMap(
				NewMapItem("c", 1),
				NewMapItem("d", 2),
				NewMapItem("e", "x"),
				NewMapItem("f", "z"),
				NewMapItem("g", LinkedHMap(
					NewMapItem("g1", 5),
				)),
			)),
			NewMapItem("a.b.c", 3),
			NewMapItem("a", LinkedHMap(
				NewMapItem("b.e", "y"),
				NewMapItem("b", LinkedHMap(
					NewMapItem("f", ""),
				)),
			)),
			NewMapItem("a.b.h", "Not NIL"),
		), ".")

		require.Equal(t, 1, actual.Size())
		a := assertChildMap(t, actual, "a")
		b := assertChildMap(t, a, "b")
		assertEquals(t, b, "c", 3)
		assertEquals(t, b, "d", 2)
		assertEquals(t, b, "e", "y")
		assertEquals(t, b, "f", "")

		g := assertChildMap(t, b, "g")
		assertEquals(t, g, "g1", 5)

		assertEquals(t, b, "h", "Not NIL")
	})
}

func Test_YamlMapSliceToLinkedHashMap(t *testing.T) {
	t.Run("Give yaml.MapSlice Should Return Success", func(t *testing.T) {
		in := yaml.MapSlice{
			{
				"a", 1,
			},
			{
				"b", yaml.MapItem{Key: "c", Value: 2},
			},
			{
				Key: "d",
				Value: yaml.MapSlice{
					{"e", 3},
					{"f", "v3"},
				},
			},
			{
				Key: "g",
				Value: []yaml.MapSlice{
					{
						yaml.MapItem{Key: "h", Value: 4},
						yaml.MapItem{Key: "i", Value: 5},
					},
					{
						yaml.MapItem{Key: "k", Value: 6},
					},
				},
			},
		}
		got := YamlMapSliceToLinkedHMap(in)
		require.Equal(t, 4, got.Size())
		v, found := got.Get("a")
		require.True(t, found)
		require.Equal(t, 1, v)

		v, found = got.Get("b")
		require.True(t, found)
		require.IsType(t, &linkedhashmap.Map{}, v)

		v1, found := v.(*linkedhashmap.Map).Get("c")
		require.True(t, found)
		require.Equal(t, 2, v1)

		v, found = got.Get("d")
		require.True(t, found)
		require.IsType(t, &linkedhashmap.Map{}, v)

		v1, found = v.(*linkedhashmap.Map).Get("e")
		require.True(t, found)
		require.Equal(t, 3, v1)

		v1, found = v.(*linkedhashmap.Map).Get("f")
		require.True(t, found)
		require.Equal(t, "v3", v1)

		v, found = got.Get("g")
		require.True(t, found)
		require.IsType(t, []*linkedhashmap.Map{}, v)
		require.Len(t, v.([]*linkedhashmap.Map), 2)

		v1 = v.([]*linkedhashmap.Map)[0]
		v2, found := v1.(*linkedhashmap.Map).Get("h")
		require.True(t, found)
		require.Equal(t, 4, v2)

		v2, found = v1.(*linkedhashmap.Map).Get("i")
		require.True(t, found)
		require.Equal(t, 5, v2)

		v1 = v.([]*linkedhashmap.Map)[1]
		v2, found = v1.(*linkedhashmap.Map).Get("k")
		require.True(t, found)
		require.Equal(t, 6, v2)
	})
}
