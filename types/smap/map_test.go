package smap_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/infraboard/mcube/types/smap"
)

type MapParameter = smap.StringMap

type testStruct struct {
	Name string
	MP   MapParameter
}

func TestMapParameterJSON(t *testing.T) {
	should := assert.New(t)
	a := MapParameter{
		"a": 1,
		"b": 2,
	}
	data, err := json.Marshal(a)
	should.NoError(err)
	should.Equal(`{"a":1,"b":2}`, string(data))
	b := testStruct{Name: "xxx", MP: a}
	data, err = json.Marshal(b)
	should.NoError(err)
	should.Equal(`{"Name":"xxx","MP":{"a":1,"b":2}}`, string(data))
}

func TestMapParameterUpdate(t *testing.T) {
	assert := assert.New(t)
	a := MapParameter{
		"a": 1,
		"b": 2,
	}
	b := MapParameter{
		"b": 3,
		"c": 4,
	}
	a.Update(b)
	assert.Equal(a, smap.StringMap{"a": 1, "b": 3, "c": 4})
}
func TestMapParameterDeepUpdate(t *testing.T) {
	tests := []struct {
		a, b, expected MapParameter
	}{
		{
			MapParameter{"a": 1},
			MapParameter{"b": 2},
			MapParameter{"a": 1, "b": 2},
		},
		{
			MapParameter{"a": 1},
			MapParameter{"a": 2},
			MapParameter{"a": 2},
		},
		{
			MapParameter{"a": 1},
			MapParameter{"a": MapParameter{"b": 1}},
			MapParameter{"a": MapParameter{"b": 1}},
		},
		{
			MapParameter{"a": MapParameter{"b": 1}},
			MapParameter{"a": MapParameter{"c": 2}},
			MapParameter{"a": MapParameter{"b": 1, "c": 2}},
		},
		{
			MapParameter{"a": MapParameter{"b": 1}},
			MapParameter{"a": 1},
			MapParameter{"a": 1},
		},
		{
			MapParameter{"a.b": 1},
			MapParameter{"a": 1},
			MapParameter{"a": 1, "a.b": 1},
		},
		{
			MapParameter{"a": 1},
			MapParameter{"a.b": 1},
			MapParameter{"a": 1, "a.b": 1},
		},
		{
			MapParameter{"a": (MapParameter)(nil)},
			MapParameter{"a": MapParameter{"b": 1}},
			MapParameter{"a": MapParameter{"b": 1}},
		},
	}
	for i, test := range tests {
		a, b, expected := test.a, test.b, test.expected
		name := fmt.Sprintf("%v: %v + %v = %v", i, a, b, expected)
		t.Run(name, func(t *testing.T) {
			a.DeepUpdate(b)
			assert.Equal(t, expected, a)
		})
	}
}
func TestMapParameterUnion(t *testing.T) {
	assert := assert.New(t)
	a := MapParameter{
		"a": 1,
		"b": 2,
	}
	b := MapParameter{
		"b": 3,
		"c": 4,
	}
	c := smap.MapStrUnion(a, b)
	assert.Equal(c, MapParameter{"a": 1, "b": 3, "c": 4})
}
func TestMapParameterCopyFieldsTo(t *testing.T) {
	assert := assert.New(t)
	m := MapParameter{
		"a": MapParameter{
			"a1": 2,
			"a2": 3,
		},
		"b": 2,
		"c": MapParameter{
			"c1": 1,
			"c2": 2,
			"c3": MapParameter{
				"c31": 1,
				"c32": 2,
			},
		},
	}
	c := MapParameter{}
	err := m.CopyFieldsTo(c, "dd")
	assert.Error(err)
	assert.Equal(MapParameter{}, c)
	err = m.CopyFieldsTo(c, "a")
	assert.Equal(nil, err)
	assert.Equal(MapParameter{"a": MapParameter{"a1": 2, "a2": 3}}, c)
	err = m.CopyFieldsTo(c, "c.c1")
	assert.Equal(nil, err)
	assert.Equal(MapParameter{"a": MapParameter{"a1": 2, "a2": 3}, "c": MapParameter{"c1": 1}}, c)
	err = m.CopyFieldsTo(c, "b")
	assert.Equal(nil, err)
	assert.Equal(MapParameter{"a": MapParameter{"a1": 2, "a2": 3}, "c": MapParameter{"c1": 1}, "b": 2}, c)
	err = m.CopyFieldsTo(c, "c.c3.c32")
	assert.Equal(nil, err)
	assert.Equal(MapParameter{"a": MapParameter{"a1": 2, "a2": 3}, "c": MapParameter{"c1": 1, "c3": MapParameter{"c32": 2}}, "b": 2}, c)
}
func TestMapParameterDelete(t *testing.T) {
	assert := assert.New(t)
	m := MapParameter{
		"c": MapParameter{
			"c1": 1,
			"c2": 2,
			"c3": MapParameter{
				"c31": 1,
				"c32": 2,
			},
		},
	}
	err := m.Delete("c.c2")
	assert.Equal(nil, err)
	assert.Equal(MapParameter{"c": MapParameter{"c1": 1, "c3": MapParameter{"c31": 1, "c32": 2}}}, m)
	err = m.Delete("c.c2.c21")
	assert.NotEqual(nil, err)
	assert.Equal(MapParameter{"c": MapParameter{"c1": 1, "c3": MapParameter{"c31": 1, "c32": 2}}}, m)
	err = m.Delete("c.c3.c31")
	assert.Equal(nil, err)
	assert.Equal(MapParameter{"c": MapParameter{"c1": 1, "c3": MapParameter{"c32": 2}}}, m)
	err = m.Delete("c")
	assert.Equal(nil, err)
	assert.Equal(MapParameter{}, m)
}
func TestHasKey(t *testing.T) {
	assert := assert.New(t)
	m := MapParameter{
		"c": MapParameter{
			"c1": 1,
			"c2": 2,
			"c3": MapParameter{
				"c31": 1,
				"c32": 2,
			},
			"c4.f": 19,
		},
		"d.f": 1,
	}
	hasKey, err := m.HasKey("c.c2")
	assert.Equal(nil, err)
	assert.Equal(true, hasKey)
	hasKey, err = m.HasKey("c.c4")
	assert.Equal(nil, err)
	assert.Equal(false, hasKey)
	hasKey, err = m.HasKey("c.c3.c32")
	assert.Equal(nil, err)
	assert.Equal(true, hasKey)
	hasKey, err = m.HasKey("dd")
	assert.Equal(nil, err)
	assert.Equal(false, hasKey)
	hasKey, err = m.HasKey("d.f")
	assert.Equal(nil, err)
	assert.Equal(true, hasKey)
	hasKey, err = m.HasKey("c.c4.f")
	assert.Equal(nil, err)
	assert.Equal(true, hasKey)
}
func TestMapParameterPut(t *testing.T) {
	m := MapParameter{
		"subMap": MapParameter{
			"a": 1,
		},
	}
	// Add new value to the top-level.
	v, err := m.Put("a", "ok")
	assert.NoError(t, err)
	assert.Nil(t, v)
	assert.Equal(t, MapParameter{"a": "ok", "subMap": MapParameter{"a": 1}}, m)
	// Add new value to subMap.
	v, err = m.Put("subMap.b", 2)
	assert.NoError(t, err)
	assert.Nil(t, v)
	assert.Equal(t, MapParameter{"a": "ok", "subMap": MapParameter{"a": 1, "b": 2}}, m)
	// Overwrite a value in subMap.
	v, err = m.Put("subMap.a", 2)
	assert.NoError(t, err)
	assert.Equal(t, 1, v)
	assert.Equal(t, MapParameter{"a": "ok", "subMap": MapParameter{"a": 2, "b": 2}}, m)
	// Add value to map that does not exist.
	m = MapParameter{}
	v, err = m.Put("subMap.newMap.a", 1)
	assert.NoError(t, err)
	assert.Nil(t, v)
	assert.Equal(t, MapParameter{"subMap": MapParameter{"newMap": MapParameter{"a": 1}}}, m)
}
func TestMapParameterGetValue(t *testing.T) {
	tests := []struct {
		input  MapParameter
		key    string
		output interface{}
		error  bool
	}{
		{
			MapParameter{"a": 1},
			"a",
			1,
			false,
		},
		{
			MapParameter{"a": MapParameter{"b": 1}},
			"a",
			MapParameter{"b": 1},
			false,
		},
		{
			MapParameter{"a": MapParameter{"b": 1}},
			"a.b",
			1,
			false,
		},
		{
			MapParameter{"a": MapParameter{"b.c": 1}},
			"a",
			MapParameter{"b.c": 1},
			false,
		},
		{
			MapParameter{"a": MapParameter{"b.c": 1}},
			"a.b",
			nil,
			true,
		},
		{
			MapParameter{"a.b": MapParameter{"c": 1}},
			"a.b",
			MapParameter{"c": 1},
			false,
		},
		{
			MapParameter{"a.b": MapParameter{"c": 1}},
			"a.b.c",
			nil,
			true,
		},
		{
			MapParameter{"a": MapParameter{"b.c": 1}},
			"a.b.c",
			1,
			false,
		},
	}
	for _, test := range tests {
		v, err := test.input.Get(test.key)
		if test.error {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, test.output, v)
	}
}
func TestClone(t *testing.T) {
	assert := assert.New(t)
	m := MapParameter{
		"c1": 1,
		"c2": 2,
		"c3": MapParameter{
			"c31": 1,
			"c32": 2,
		},
	}
	c := m.Clone()
	assert.Equal(MapParameter{"c31": 1, "c32": 2}, c["c3"])
}
func TestString(t *testing.T) {
	type io struct {
		Input  MapParameter
		Output string
	}
	tests := []io{
		{
			Input: MapParameter{
				"a": "b",
			},
			Output: `{"a":"b"}`,
		},
		{
			Input: MapParameter{
				"a": []int{1, 2, 3},
			},
			Output: `{"a":[1,2,3]}`,
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.Output, test.Input.String())
	}
}

// Smoke test. The method has no observable outputs so this
// is only verifying there are no panics.
func TestStringToPrint(t *testing.T) {
	m := MapParameter{}
	assert.Equal(t, "{}", m.StringToPrint())
	assert.Equal(t, true, len(m.StringToPrint()) > 0)
}
func TestMergeFields(t *testing.T) {
	type io struct {
		UnderRoot bool
		Event     MapParameter
		Fields    MapParameter
		Output    MapParameter
		Err       string
	}
	tests := []io{
		// underRoot = true, merges
		{
			UnderRoot: true,
			Event: MapParameter{
				"a": "1",
			},
			Fields: MapParameter{
				"b": 2,
			},
			Output: MapParameter{
				"a": "1",
				"b": 2,
			},
		},
		// underRoot = true, overwrites existing
		{
			UnderRoot: true,
			Event: MapParameter{
				"a": "1",
			},
			Fields: MapParameter{
				"a": 2,
			},
			Output: MapParameter{
				"a": 2,
			},
		},
		// underRoot = false, adds new 'fields' when it doesn't exist
		{
			UnderRoot: false,
			Event: MapParameter{
				"a": "1",
			},
			Fields: MapParameter{
				"a": 2,
			},
			Output: MapParameter{
				"a": "1",
				"fields": MapParameter{
					"a": 2,
				},
			},
		},
		// underRoot = false, merge with existing 'fields' and overwrites existing keys
		{
			UnderRoot: false,
			Event: MapParameter{
				"fields": MapParameter{
					"a": "1",
					"b": 2,
				},
			},
			Fields: MapParameter{
				"a": 3,
				"c": 4,
			},
			Output: MapParameter{
				"fields": MapParameter{
					"a": 3,
					"b": 2,
					"c": 4,
				},
			},
		},
		// underRoot = false, error when 'fields' is wrong type
		{
			UnderRoot: false,
			Event: MapParameter{
				"fields": "not a MapParameter",
			},
			Fields: MapParameter{
				"a": 3,
			},
			Output: MapParameter{
				"fields": "not a MapParameter",
			},
			Err: "expected map",
		},
	}
	for _, test := range tests {
		err := smap.MergeFields(test.Event, test.Fields, test.UnderRoot)
		assert.Equal(t, test.Output, test.Event)
		if test.Err != "" {
			assert.Contains(t, err.Error(), test.Err)
		} else {
			assert.NoError(t, err)
		}
	}
}
func TestMergeFieldsDeep(t *testing.T) {
	type io struct {
		UnderRoot bool
		Event     MapParameter
		Fields    MapParameter
		Output    MapParameter
		Err       string
	}
	tests := []io{
		// underRoot = true, merges
		{
			UnderRoot: true,
			Event: MapParameter{
				"a": "1",
			},
			Fields: MapParameter{
				"b": 2,
			},
			Output: MapParameter{
				"a": "1",
				"b": 2,
			},
		},
		// underRoot = true, overwrites existing
		{
			UnderRoot: true,
			Event: MapParameter{
				"a": "1",
			},
			Fields: MapParameter{
				"a": 2,
			},
			Output: MapParameter{
				"a": 2,
			},
		},
		// underRoot = false, adds new 'fields' when it doesn't exist
		{
			UnderRoot: false,
			Event: MapParameter{
				"a": "1",
			},
			Fields: MapParameter{
				"a": 2,
			},
			Output: MapParameter{
				"a": "1",
				"fields": MapParameter{
					"a": 2,
				},
			},
		},
		// underRoot = false, merge with existing 'fields' and overwrites existing keys
		{
			UnderRoot: false,
			Event: MapParameter{
				"fields": MapParameter{
					"a": "1",
					"b": 2,
				},
			},
			Fields: MapParameter{
				"a": 3,
				"c": 4,
			},
			Output: MapParameter{
				"fields": MapParameter{
					"a": 3,
					"b": 2,
					"c": 4,
				},
			},
		},
		// underRoot = false, error when 'fields' is wrong type
		{
			UnderRoot: false,
			Event: MapParameter{
				"fields": "not a MapParameter",
			},
			Fields: MapParameter{
				"a": 3,
			},
			Output: MapParameter{
				"fields": "not a MapParameter",
			},
			Err: "expected map",
		},
		// underRoot = true, merges recursively
		{
			UnderRoot: true,
			Event: MapParameter{
				"my": MapParameter{
					"field1": "field1",
				},
			},
			Fields: MapParameter{
				"my": MapParameter{
					"field2": "field2",
					"field3": "field3",
				},
			},
			Output: MapParameter{
				"my": MapParameter{
					"field1": "field1",
					"field2": "field2",
					"field3": "field3",
				},
			},
		},
		// underRoot = true, merges recursively and overrides
		{
			UnderRoot: true,
			Event: MapParameter{
				"my": MapParameter{
					"field1": "field1",
					"field2": "field2",
				},
			},
			Fields: MapParameter{
				"my": MapParameter{
					"field2": "fieldTWO",
					"field3": "field3",
				},
			},
			Output: MapParameter{
				"my": MapParameter{
					"field1": "field1",
					"field2": "fieldTWO",
					"field3": "field3",
				},
			},
		},
		// underRoot = false, merges recursively under existing 'fields'
		{
			UnderRoot: false,
			Event: MapParameter{
				"fields": MapParameter{
					"my": MapParameter{
						"field1": "field1",
					},
				},
			},
			Fields: MapParameter{
				"my": MapParameter{
					"field2": "field2",
					"field3": "field3",
				},
			},
			Output: MapParameter{
				"fields": MapParameter{
					"my": MapParameter{
						"field1": "field1",
						"field2": "field2",
						"field3": "field3",
					},
				},
			},
		},
	}
	for _, test := range tests {
		err := smap.MergeFieldsDeep(test.Event, test.Fields, test.UnderRoot)
		assert.Equal(t, test.Output, test.Event)
		if test.Err != "" {
			assert.Contains(t, err.Error(), test.Err)
		} else {
			assert.NoError(t, err)
		}
	}
}
func TestAddTag(t *testing.T) {
	type io struct {
		Event  MapParameter
		Tags   []string
		Output MapParameter
		Err    string
	}
	tests := []io{
		// No existing tags, creates new tag array
		{
			Event: MapParameter{},
			Tags:  []string{"json"},
			Output: MapParameter{
				"tags": []string{"json"},
			},
		},
		// Existing tags is a []string, appends
		{
			Event: MapParameter{
				"tags": []string{"json"},
			},
			Tags: []string{"docker"},
			Output: MapParameter{
				"tags": []string{"json", "docker"},
			},
		},
		// Existing tags is a []interface{}, appends
		{
			Event: MapParameter{
				"tags": []interface{}{"json"},
			},
			Tags: []string{"docker"},
			Output: MapParameter{
				"tags": []interface{}{"json", "docker"},
			},
		},
		// Existing tags is not a []string or []interface{}
		{
			Event: MapParameter{
				"tags": "not a slice",
			},
			Tags: []string{"docker"},
			Output: MapParameter{
				"tags": "not a slice",
			},
			Err: "expected string array",
		},
	}
	for _, test := range tests {
		err := smap.AddTags(test.Event, test.Tags)
		assert.Equal(t, test.Output, test.Event)
		if test.Err != "" {
			assert.Contains(t, err.Error(), test.Err)
		} else {
			assert.NoError(t, err)
		}
	}
}
func TestAddTagsWithKey(t *testing.T) {
	type io struct {
		Event  MapParameter
		Key    string
		Tags   []string
		Output MapParameter
		Err    string
	}
	tests := []io{
		// No existing tags, creates new tag array
		{
			Event: MapParameter{},
			Key:   "tags",
			Tags:  []string{"json"},
			Output: MapParameter{
				"tags": []string{"json"},
			},
		},
		// Existing tags is a []string, appends
		{
			Event: MapParameter{
				"tags": []string{"json"},
			},
			Key:  "tags",
			Tags: []string{"docker"},
			Output: MapParameter{
				"tags": []string{"json", "docker"},
			},
		},
		// Existing tags are in submap and is a []interface{}, appends
		{
			Event: MapParameter{
				"log": MapParameter{
					"flags": []interface{}{"json"},
				},
			},
			Key:  "log.flags",
			Tags: []string{"docker"},
			Output: MapParameter{
				"log": MapParameter{
					"flags": []interface{}{"json", "docker"},
				},
			},
		},
		// Existing tags are in a submap and is not a []string or []interface{}
		{
			Event: MapParameter{
				"log": MapParameter{
					"flags": "not a slice",
				},
			},
			Key:  "log.flags",
			Tags: []string{"docker"},
			Output: MapParameter{
				"log": MapParameter{
					"flags": "not a slice",
				},
			},
			Err: "expected string array",
		},
	}
	for _, test := range tests {
		err := smap.AddTagsWithKey(test.Event, test.Key, test.Tags)
		assert.Equal(t, test.Output, test.Event)
		if test.Err != "" {
			assert.Contains(t, err.Error(), test.Err)
		} else {
			assert.NoError(t, err)
		}
	}
}
func TestFlatten(t *testing.T) {
	type data struct {
		Event    MapParameter
		Expected MapParameter
	}
	tests := []data{
		{
			Event: MapParameter{
				"hello": MapParameter{
					"world": 15,
				},
			},
			Expected: MapParameter{
				"hello.world": 15,
			},
		},
		{
			Event: MapParameter{
				"test": 15,
			},
			Expected: MapParameter{
				"test": 15,
			},
		},
		{
			Event: MapParameter{
				"test": 15,
				"hello": MapParameter{
					"world": MapParameter{
						"ok": "test",
					},
				},
				"elastic": MapParameter{
					"for": "search",
				},
			},
			Expected: MapParameter{
				"test":           15,
				"hello.world.ok": "test",
				"elastic.for":    "search",
			},
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.Expected, test.Event.Flatten())
	}
}
func BenchmarkMapParameterFlatten(b *testing.B) {
	m := MapParameter{
		"test": 15,
		"hello": MapParameter{
			"world": MapParameter{
				"ok": "test",
			},
		},
		"elastic": MapParameter{
			"for": "search",
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Flatten()
	}
}
func BenchmarkWalkMap(b *testing.B) {
	globalM := MapParameter{
		"hello": MapParameter{
			"world": MapParameter{
				"ok": "test",
			},
		},
	}
	b.Run("Get", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			globalM.Get("test.world.ok")
		}
	})
	b.Run("Put", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			m := MapParameter{
				"hello": MapParameter{
					"world": MapParameter{
						"ok": "test",
					},
				},
			}
			m.Put("hello.world.new", 17)
		}
	})
	b.Run("PutMissing", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			m := MapParameter{}
			m.Put("a.b.c", 17)
		}
	})
	b.Run("HasKey", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			globalM.HasKey("hello.world.ok")
			globalM.HasKey("hello.world.no_ok")
		}
	})
	b.Run("HasKeyFirst", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			globalM.HasKey("hello")
		}
	})
	b.Run("Delete", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			m := MapParameter{
				"hello": MapParameter{
					"world": MapParameter{
						"ok": "test",
					},
				},
			}
			m.Put("hello.world.test", 17)
		}
	})
}
