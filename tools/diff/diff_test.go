package diff_test

import (
	"errors"
	"testing"
	"time"

	"github.com/infraboard/mcube/v2/tools/diff"
)

type item struct {
	ID   string `diff:"ID,key"`
	Name string `diff:"名称"`
}

type root struct {
	Items []*item         `diff:"列表"`
	Meta  map[string]item `diff:"元数据"`
	Note  any             `diff:"备注"`
	Tags  []string        `diff:"标签"`
	Name  string          `diff:"名称,level=info"`
}

func TestCompareE_NilAny(t *testing.T) {
	recs, err := diff.CompareE(nil, nil)
	if err != nil || len(recs) != 0 {
		t.Fatalf("both nil: recs=%v err=%v", recs, err)
	}
	_, err = diff.CompareE(nil, "x")
	if !errors.Is(err, diff.ErrTypeMismatch) {
		t.Fatalf("nil vs value: err=%v", err)
	}
}

func TestCompareE_TypeMismatch(t *testing.T) {
	_, err := diff.CompareE(1, "1")
	if !errors.Is(err, diff.ErrTypeMismatch) {
		t.Fatalf("err=%v", err)
	}
}

func TestCompare_BasicAndTime(t *testing.T) {
	type s struct {
		N int       `diff:"数量"`
		T time.Time `diff:"时间"`
	}
	t1 := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	t2 := t1.Add(time.Hour)
	recs := diff.Compare(s{N: 1, T: t1}, s{N: 2, T: t2})
	if len(recs) != 2 {
		t.Fatalf("want 2 diffs, got %#v", recs)
	}
}

func TestCompare_InterfaceField(t *testing.T) {
	a := root{Note: "old"}
	b := root{Note: "new"}
	recs := diff.Compare(a, b)
	if len(recs) != 1 || recs[0].FieldPath != "Note" {
		t.Fatalf("want Note change, got %#v", recs)
	}
	if recs[0].OldValue != "old" || recs[0].NewValue != "new" {
		t.Fatalf("values: %#v", recs[0])
	}
}

func TestCompare_SliceInPlacePathIncludesKey(t *testing.T) {
	type s struct {
		Items []item `diff:"列表"`
	}
	a := s{Items: []item{{ID: "1", Name: "a"}, {ID: "2", Name: "b"}}}
	b := s{Items: []item{{ID: "1", Name: "A"}, {ID: "2", Name: "B"}}}
	recs := diff.Compare(a, b)
	if len(recs) != 2 {
		t.Fatalf("want 2, got %#v", recs)
	}
	paths := map[string]bool{}
	for _, r := range recs {
		paths[r.FieldPath] = true
	}
	if !paths["Items[1].Name"] || !paths["Items[2].Name"] {
		t.Fatalf("paths=%v", paths)
	}
}

func TestCompare_PointerSliceAdd(t *testing.T) {
	a := root{Items: []*item{{ID: "1", Name: "a"}}}
	b := root{Items: []*item{{ID: "1", Name: "a"}, {ID: "2", Name: "b"}}}
	recs := diff.CompareNormalized(a, b)
	found := false
	for _, r := range recs {
		if r.FieldPath == "Items.ID" && r.OldValue == diff.MarkerAdd && r.NewValue == "2" {
			found = true
		}
	}
	if !found {
		t.Fatalf("want tagged add for []*item, got %#v", recs)
	}
}

func TestCompare_EmptyStringKeyAdd(t *testing.T) {
	type e struct {
		ID string `diff:"ID,key"`
	}
	type s struct {
		Items []e `diff:"列表"`
	}
	recs := diff.Compare(s{}, s{Items: []e{{ID: ""}}})
	if len(recs) != 1 {
		t.Fatalf("want 1 add, got %#v", recs)
	}
	if recs[0].OldValue != diff.MarkerAdd || recs[0].NewValue != "" {
		t.Fatalf("record=%#v", recs[0])
	}
}

func TestCompare_MapValueFieldDiff(t *testing.T) {
	a := root{Meta: map[string]item{"x": {ID: "1", Name: "a"}}}
	b := root{Meta: map[string]item{"x": {ID: "1", Name: "b"}}}
	recs := diff.Compare(a, b)
	if len(recs) != 1 {
		t.Fatalf("want 1 field diff, got %#v", recs)
	}
	if recs[0].FieldPath != "Meta[x].Name" || recs[0].OldValue != "a" || recs[0].NewValue != "b" {
		t.Fatalf("record=%#v", recs[0])
	}
}

func TestCompare_FieldLevelsRespected(t *testing.T) {
	type s struct {
		Name string `diff:"名称,level=info"`
	}
	opt := &diff.Options{FieldLevels: map[string]diff.DiffLevel{"Name": diff.LevelError}}
	recs := diff.Compare(s{Name: "a"}, s{Name: "b"}, opt)
	if len(recs) != 1 || recs[0].Level != diff.LevelError {
		t.Fatalf("want LevelError, got %#v", recs)
	}
	// 不应污染调用方 map
	if len(opt.FieldLevels) != 1 {
		t.Fatalf("options mutated: %#v", opt.FieldLevels)
	}
}

func TestCompare_IgnoreAndOnlyTagged(t *testing.T) {
	type s struct {
		A string `diff:"A"`
		B string `diff:"-"`
		C string
	}
	recs := diff.Compare(
		s{A: "1", B: "1", C: "1"},
		s{A: "2", B: "2", C: "2"},
		&diff.Options{OnlyTaggedFields: true},
	)
	if len(recs) != 1 || recs[0].FieldPath != "A" {
		t.Fatalf("got %#v", recs)
	}
}

func TestNormalizeNilContainers_DeepCopyAndNilSlice(t *testing.T) {
	type inner struct {
		Tags []string
	}
	type s struct {
		M map[string]inner
		L []string
	}
	orig := s{M: map[string]inner{"k": {Tags: nil}}}
	cp := diff.NormalizeNilContainers(orig)
	if cp.L == nil || cp.M["k"].Tags == nil {
		t.Fatalf("nil containers not normalized: %#v", cp)
	}
	cp.M["k"] = inner{Tags: []string{"x"}}
	if len(orig.M["k"].Tags) != 0 {
		t.Fatalf("shallow share mutated original: %#v", orig)
	}
}

func TestCompareNormalized_NilVsEmpty(t *testing.T) {
	type s struct {
		Tags []string `diff:"标签"`
	}
	var a s
	b := s{Tags: []string{}}
	if recs := diff.CompareNormalized(a, b); len(recs) != 0 {
		t.Fatalf("nil vs empty should be equal after normalize, got %#v", recs)
	}
}

func TestMaskedChangeAndSummary(t *testing.T) {
	recs := diff.AppendRecords(nil, diff.MaskedChange("Secret", "密钥", "******", diff.LevelWarn))
	sum := diff.Summary(recs)
	if sum == "" || recs[0].Level != diff.LevelWarn {
		t.Fatalf("sum=%q recs=%#v", sum, recs)
	}
}

func TestCompare_SliceRemoveAndReorderByKey(t *testing.T) {
	type s struct {
		Items []item `diff:"列表"`
	}
	a := s{Items: []item{{ID: "1", Name: "a"}, {ID: "2", Name: "b"}}}
	b := s{Items: []item{{ID: "2", Name: "b"}, {ID: "1", Name: "a"}}}
	if recs := diff.Compare(a, b); len(recs) != 0 {
		t.Fatalf("reorder by key should be equal, got %#v", recs)
	}

	b = s{Items: []item{{ID: "2", Name: "b"}}}
	recs := diff.Compare(a, b)
	found := false
	for _, r := range recs {
		if r.NewValue == diff.MarkerRemove && r.OldValue == "1" {
			found = true
		}
	}
	if !found {
		t.Fatalf("want remove id=1, got %#v", recs)
	}
}
