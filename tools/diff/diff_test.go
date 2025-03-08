package diff_test

import (
	"testing"
	"time"

	"github.com/infraboard/mcube/v2/tools/diff"
	"github.com/stretchr/testify/assert"
)

// 测试基础结构体
type BasicStruct struct {
	ID   int    `diff:"ID"`
	Name string `diff:"名称"`
}

// 测试嵌套结构体
type NestedStruct struct {
	Basic       BasicStruct `diff:"基础信息"`
	Score       float64     `diff:"分数"`
	Secret      *string     `diff:"密钥,level=error"`
	ScoreIgnore float64     `diff:"-"`
}

// 测试切片和Map
type ComplexStruct struct {
	Tags  []string       `diff:"标签"`
	Items map[string]int `diff:"物品"`
	Time  time.Time      `diff:"时间"`
}

// 测试基础类型对比
func TestBasicTypes(t *testing.T) {
	t.Run("相同结构体", func(t *testing.T) {
		a := BasicStruct{ID: 1, Name: "Test"}
		b := BasicStruct{ID: 1, Name: "Test"}
		diffs := diff.Compare(a, b)
		assert.Empty(t, diffs)
	})

	t.Run("不同字段值", func(t *testing.T) {
		a := BasicStruct{ID: 1, Name: "Old"}
		b := BasicStruct{ID: 2, Name: "New"}
		diffs := diff.Compare(a, b)
		assert.Len(t, diffs, 2)
		assert.Contains(t, diffs[0].FieldPath, "ID")
		assert.Contains(t, diffs[1].FieldDesc, "名称")
	})
}

// 测试指针类型
func TestPointerTypes(t *testing.T) {
	secret1 := "key1"
	secret2 := "key2"

	tests := []struct {
		name     string
		a        NestedStruct
		b        NestedStruct
		expected int
	}{
		{
			"指针nil与非nil",
			NestedStruct{Secret: nil},
			NestedStruct{Secret: &secret1},
			1,
		},
		{
			"指针值不同",
			NestedStruct{Secret: &secret1},
			NestedStruct{Secret: &secret2},
			1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diffs := diff.Compare(tt.a, tt.b)
			assert.Len(t, diffs, tt.expected)
		})
	}
}

// 测试嵌套结构体
func TestNestedStruct(t *testing.T) {
	a := NestedStruct{
		Basic:  BasicStruct{ID: 1, Name: "Alice"},
		Score:  90.5,
		Secret: nil,
	}

	b := NestedStruct{
		Basic:  BasicStruct{ID: 2, Name: "Alice"},
		Score:  95.0,
		Secret: nil,
	}

	diffs := diff.Compare(a, b)
	assert.Len(t, diffs, 2)

	assert.Contains(t, diffs[0].FieldDesc, "基础信息.ID")
	assert.Equal(t, 1, diffs[0].OldValue)
	assert.Equal(t, 2, diffs[0].NewValue)
	assert.Equal(t, 90.5, diffs[1].OldValue)
	assert.Equal(t, 95.0, diffs[1].NewValue)
}

// 测试切片/数组
func TestSliceAndArray(t *testing.T) {
	t.Run("元素删除", func(t *testing.T) {
		a := ComplexStruct{Tags: []string{"A", "B"}}
		b := ComplexStruct{Tags: []string{"A"}}
		diffs := diff.Compare(a, b)
		assert.Len(t, diffs, 1)
		assert.Contains(t, "标签[1]", diffs[0].FieldDesc)
		assert.Equal(t, "B", diffs[0].OldValue)
	})

	t.Run("元素增加", func(t *testing.T) {
		a := ComplexStruct{Tags: []string{"A"}}
		b := ComplexStruct{Tags: []string{"A", "B"}}
		diffs := diff.Compare(a, b)
		assert.Len(t, diffs, 1)
		assert.Contains(t, "标签[1]", diffs[0].FieldDesc)
		assert.Equal(t, nil, diffs[0].OldValue)
		assert.Equal(t, "B", diffs[0].NewValue)
	})

	t.Run("元素修改", func(t *testing.T) {
		a := ComplexStruct{Tags: []string{"A", "B"}}
		b := ComplexStruct{Tags: []string{"A", "C"}}
		diffs := diff.Compare(a, b)
		assert.Len(t, diffs, 1)

		assert.Contains(t, "标签[1]", diffs[0].FieldDesc)
	})
}

// 测试Map类型
func TestMapType(t *testing.T) {
	a := ComplexStruct{
		Items: map[string]int{"apple": 3, "banana": 5},
	}

	b := ComplexStruct{
		Items: map[string]int{"apple": 5, "orange": 2},
	}

	diffs := diff.Compare(a, b)
	assert.Len(t, diffs, 3) // apple值变化 + banana删除 + orange新增
}

// 测试时间类型
func TestTimeType(t *testing.T) {
	t1 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)

	a := ComplexStruct{Time: t1}
	b := ComplexStruct{Time: t2}
	diffs := diff.Compare(a, b)
	assert.Len(t, diffs, 1)
	assert.IsType(t, time.Time{}, diffs[0].OldValue)
}

// 测试配置选项
func TestIgnoreField(t *testing.T) {
	a := NestedStruct{
		Basic:       BasicStruct{ID: 1, Name: "Alice"},
		ScoreIgnore: 90.5,
		Secret:      nil,
	}

	b := NestedStruct{
		Basic:       BasicStruct{ID: 2, Name: "Alice"},
		ScoreIgnore: 95.0,
		Secret:      nil,
	}

	diffs := diff.Compare(a, b)
	assert.Len(t, diffs, 1) // 只有 Basic.ID 变化
	assert.Contains(t, diffs[0].FieldPath, "Basic.ID")
	assert.NotContains(t, diffs[0].FieldPath, "Score") // Score 字段被忽略
}

// 测试类型安全
func TestTypeSafety(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("未触发类型不匹配panic")
		}
	}()

	a := BasicStruct{ID: 1}
	b := struct{ ID string }{ID: "1"}
	diff.Compare(a, b)
}

// 测试性能（可选）
func BenchmarkCompare(b *testing.B) {
	type deepStruct struct {
		Level1 struct {
			Level2 struct {
				Value int `diff:"值"`
			} `diff:"二层"`
		} `diff:"一层"`
	}

	a := deepStruct{}
	a.Level1.Level2.Value = 1

	bigData := ComplexStruct{
		Tags:  make([]string, 1000),
		Items: make(map[string]int, 1000),
	}

	b.Run("深层结构", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			diff.Compare(a, a)
		}
	})

	b.Run("大数据量", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			diff.Compare(bigData, bigData)
		}
	})
}

type ComplexNestedStruct struct {
	Basic       []BasicStruct `diff:"基础信息"`
	Score       float64       `diff:"分数"`
	Secret      *string       `diff:"密钥,level=error"`
	ScoreIgnore float64       `diff:"-"`
}

// 测试嵌套数组对象的添加、修改、删除
func TestNestedArrayFieldChanges(t *testing.T) {
	secretOld := "old_key"
	secretNew := "new_key"

	a := ComplexNestedStruct{
		Basic: []BasicStruct{
			{ID: 1, Name: "Alice"},
			{ID: 2, Name: "Bob"},
		},
		Score:  90.5,
		Secret: &secretOld,
	}

	b := ComplexNestedStruct{
		Basic: []BasicStruct{
			{ID: 1, Name: "Alice"},
			{ID: 3, Name: "Charlie"},
		},
		Score:  95.0,
		Secret: &secretNew,
	}

	diffs := diff.Compare(a, b)

	expected := []diff.DiffRecord{
		{
			FieldPath: "Basic[1].ID",
			FieldDesc: "基础信息.ID",
			OldValue:  2,
			NewValue:  3,
			Level:     diff.LevelInfo,
		},
		{
			FieldPath: "Basic[1].Name",
			FieldDesc: "基础信息.名称",
			OldValue:  "Bob",
			NewValue:  "Charlie",
			Level:     diff.LevelInfo,
		},
		{
			FieldPath: "Score",
			FieldDesc: "分数",
			OldValue:  90.5,
			NewValue:  95.0,
			Level:     diff.LevelInfo,
		},
		{
			FieldPath: "Secret",
			FieldDesc: "密钥",
			OldValue:  "old_key",
			NewValue:  "new_key",
			Level:     diff.LevelError,
		},
	}

	assert.Equal(t, len(expected), len(diffs), "差异记录数量不匹配")
	for i, exp := range expected {
		assert.Equal(t, exp.FieldPath, diffs[i].FieldPath, "字段路径不匹配")
		assert.Equal(t, exp.FieldDesc, diffs[i].FieldDesc, "字段描述不匹配")
		assert.Equal(t, exp.OldValue, diffs[i].OldValue, "旧值不匹配")
		assert.Equal(t, exp.NewValue, diffs[i].NewValue, "新值不匹配")
		assert.Equal(t, exp.Level, diffs[i].Level, "差异级别不匹配")
	}
}

func TestNestedArrayFieldAdd(t *testing.T) {
	a := ComplexNestedStruct{
		Basic: []BasicStruct{
			{ID: 1, Name: "Alice"},
		},
	}

	b := ComplexNestedStruct{
		Basic: []BasicStruct{
			{ID: 1, Name: "Alice"},
			{ID: 3, Name: "Charlie"},
		},
	}

	diffs := diff.Compare(a, b)

	expected := []diff.DiffRecord{
		{
			FieldPath: "Basic[1].ID",
			FieldDesc: "基础信息.ID",
			OldValue:  nil,
			NewValue:  3,
			Level:     diff.LevelInfo,
		},
		{
			FieldPath: "Basic[1].Name",
			FieldDesc: "基础信息.名称",
			OldValue:  nil,
			NewValue:  "Charlie",
			Level:     diff.LevelInfo,
		},
	}

	assert.Equal(t, len(expected), len(diffs), "差异记录数量不匹配")
	for i, exp := range expected {
		assert.Equal(t, exp.FieldPath, diffs[i].FieldPath, "字段路径不匹配")
		assert.Equal(t, exp.FieldDesc, diffs[i].FieldDesc, "字段描述不匹配")
		assert.Equal(t, exp.OldValue, diffs[i].OldValue, "旧值不匹配")
		assert.Equal(t, exp.NewValue, diffs[i].NewValue, "新值不匹配")
		assert.Equal(t, exp.Level, diffs[i].Level, "差异级别不匹配")
	}
}

func TestNestedArrayFieldRemove(t *testing.T) {
	a := ComplexNestedStruct{
		Basic: []BasicStruct{
			{ID: 1, Name: "Alice"},
			{ID: 3, Name: "Charlie"},
		},
	}

	b := ComplexNestedStruct{
		Basic: []BasicStruct{
			{ID: 1, Name: "Alice"},
		},
	}

	diffs := diff.Compare(a, b)

	expected := []diff.DiffRecord{
		{
			FieldPath: "Basic[1].ID",
			FieldDesc: "基础信息.ID",
			OldValue:  3,
			NewValue:  nil,
			Level:     diff.LevelInfo,
		},
		{
			FieldPath: "Basic[1].Name",
			FieldDesc: "基础信息.名称",
			OldValue:  "Charlie",
			NewValue:  nil,
			Level:     diff.LevelInfo,
		},
	}

	assert.Equal(t, len(expected), len(diffs), "差异记录数量不匹配")
	for i, exp := range expected {
		assert.Equal(t, exp.FieldPath, diffs[i].FieldPath, "字段路径不匹配")
		assert.Equal(t, exp.FieldDesc, diffs[i].FieldDesc, "字段描述不匹配")
		assert.Equal(t, exp.OldValue, diffs[i].OldValue, "旧值不匹配")
		assert.Equal(t, exp.NewValue, diffs[i].NewValue, "新值不匹配")
		assert.Equal(t, exp.Level, diffs[i].Level, "差异级别不匹配")
	}
}
