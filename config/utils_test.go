package config

import (
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/stretchr/testify/assert"
	"gitlab.id.vin/vincart/golib/utils"
	"gopkg.in/yaml.v2"
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

func Test_expandInlineKeyInMap(t *testing.T) {
	assertExist := func(t *testing.T, m *linkedhashmap.Map, key string) interface{} {
		v, f := m.Get(key)
		assert.True(t, f)
		return v
	}

	assertChildMap := func(t *testing.T, m *linkedhashmap.Map, key string) *linkedhashmap.Map {
		v := assertExist(t, m, key)
		assert.IsType(t, &linkedhashmap.Map{}, v)
		return v.(*linkedhashmap.Map)
	}

	assertEquals := func(t *testing.T, m *linkedhashmap.Map, key string, val interface{}) {
		v := assertExist(t, m, key)
		assert.Equal(t, val, v)
	}

	t.Run("Give inline key should expand success", func(t *testing.T) {
		got := expandInlineKeyInMap(utils.LinkedHMap(
			utils.NewMapItem("a.b", utils.LinkedHMap(utils.NewMapItem("c", 1))),
			utils.NewMapItem("x", utils.LinkedHMap(utils.NewMapItem("y", []int{2, 3}))),
		), ".")
		assert.Equal(t, 2, got.Size())
		v := assertChildMap(t, got, "a")
		v1 := assertChildMap(t, v, "b")
		v2 := assertExist(t, v1, "c")
		assert.EqualValues(t, 1, v2)

		v = assertChildMap(t, got, "x")
		v3 := assertExist(t, v, "y")
		assert.Equal(t, []int{2, 3}, v3)
	})

	t.Run("Give multiple inline keys should expand success", func(t *testing.T) {
		got := expandInlineKeyInMap(utils.LinkedHMap(
			utils.NewMapItem("a.b.c", 1),
			utils.NewMapItem("x.y", []int{2, 3}),
		), ".")

		assert.Equal(t, 2, got.Size())
		v := assertChildMap(t, got, "a")
		v1 := assertChildMap(t, v, "b")
		v2 := assertExist(t, v1, "c")
		assert.EqualValues(t, 1, v2)

		v = assertChildMap(t, got, "x")
		v2 = assertExist(t, v, "y")
		assert.Equal(t, []int{2, 3}, v2)
	})

	t.Run("Give a override key should return success", func(t *testing.T) {
		actual := expandInlineKeyInMap(utils.LinkedHMap(
			utils.NewMapItem("a.b", utils.LinkedHMap(
				utils.NewMapItem("c", 1),
				utils.NewMapItem("d", 2),
				utils.NewMapItem("e", "x"),
				utils.NewMapItem("f", "z"),
				utils.NewMapItem("g", utils.LinkedHMap(
					utils.NewMapItem("g1", 5),
				)),
			)),
			utils.NewMapItem("a.b.c", 3),
			utils.NewMapItem("a", utils.LinkedHMap(
				utils.NewMapItem("b.e", "y"),
				utils.NewMapItem("b", utils.LinkedHMap(
					utils.NewMapItem("f", ""),
				)),
			)),
			utils.NewMapItem("a.b.h", "Not NIL"),
		), ".")

		assert.Equal(t, 1, actual.Size())
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

func Test_convertYamlMapSliceToLinkedHashMap(t *testing.T) {
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
		got := yamlMapSliceToLinkedHMap(in)
		assert.Equal(t, 4, got.Size())
		v, found := got.Get("a")
		assert.True(t, found)
		assert.Equal(t, 1, v)

		v, found = got.Get("b")
		assert.True(t, found)
		assert.IsType(t, &linkedhashmap.Map{}, v)

		v1, found := v.(*linkedhashmap.Map).Get("c")
		assert.True(t, found)
		assert.Equal(t, 2, v1)

		v, found = got.Get("d")
		assert.True(t, found)
		assert.IsType(t, &linkedhashmap.Map{}, v)

		v1, found = v.(*linkedhashmap.Map).Get("e")
		assert.True(t, found)
		assert.Equal(t, 3, v1)

		v1, found = v.(*linkedhashmap.Map).Get("f")
		assert.True(t, found)
		assert.Equal(t, "v3", v1)

		v, found = got.Get("g")
		assert.True(t, found)
		assert.IsType(t, []*linkedhashmap.Map{}, v)
		assert.Len(t, v.([]*linkedhashmap.Map), 2)

		v1 = v.([]*linkedhashmap.Map)[0]
		v2, found := v1.(*linkedhashmap.Map).Get("h")
		assert.True(t, found)
		assert.Equal(t, 4, v2)

		v2, found = v1.(*linkedhashmap.Map).Get("i")
		assert.True(t, found)
		assert.Equal(t, 5, v2)

		v1 = v.([]*linkedhashmap.Map)[1]
		v2, found = v1.(*linkedhashmap.Map).Get("k")
		assert.True(t, found)
		assert.Equal(t, 6, v2)
	})
}

func Test_mergeLinkedHashMap(t *testing.T) {
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

		mergeLinkedHMap(dst, src)
		assert.Equal(t, 4, dst.Size())

		v, found := dst.Get("a")
		assert.True(t, found)
		assert.EqualValues(t, 1, v)

		v, found = dst.Get("b")
		assert.True(t, found)
		assert.EqualValues(t, "v2-modified", v)

		v, found = dst.Get("c")
		assert.True(t, found)
		assert.IsType(t, &linkedhashmap.Map{}, v)

		vc, found := v.(*linkedhashmap.Map).Get("d")
		assert.True(t, found)
		assert.EqualValues(t, 3, vc)

		vc, found = v.(*linkedhashmap.Map).Get("e")
		assert.True(t, found)
		assert.EqualValues(t, 5, vc)

		vc, found = v.(*linkedhashmap.Map).Get("f")
		assert.True(t, found)
		assert.EqualValues(t, "v6", vc)

		v, found = dst.Get("g")
		assert.True(t, found)
		assert.EqualValues(t, 7, v)
	})
}
