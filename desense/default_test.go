package desense_test

import (
	"testing"

	"github.com/infraboard/mcube/v2/desense"
	"github.com/infraboard/mcube/v2/tools/pretty"
)

func TestDeSense(t *testing.T) {
	v := desense.Default().DeSense("abcdefghijkl", "3", "2")
	t.Logf(v)
}

type TestMaskStuct struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number" mask:",5,4"`
	TestEmbedMaskStruct
}

type TestEmbedMaskStruct struct {
	Name1        string `json:"name1"`
	PhoneNumber1 string `json:"phone_number1" mask:",5,4"`
}

func TestMusk(t *testing.T) {
	tm := &TestMaskStuct{
		PhoneNumber: "1234567890123",
		Name:        "test",
		TestEmbedMaskStruct: TestEmbedMaskStruct{
			PhoneNumber1: "1234567890123",
			Name1:        "test",
		},
	}

	if err := desense.MaskStruct(tm); err != nil {
		t.Fatal(err)
	}
	t.Log(tm.PhoneNumber1)
}

var set = &TestMaskStuctSet{
	Items: []*TestMaskStuct{
		{
			PhoneNumber: "1234567890123",
			Name:        "test",
		},
	},
}

type TestMaskStuctSet struct {
	T     TestMaskStuct    `json:"t"`
	Items []*TestMaskStuct `json:"items"`
}

func TestMuskSet(t *testing.T) {
	if err := desense.MaskStruct(set); err != nil {
		t.Fatal(err)
	}

	t.Log(pretty.ToJSON(set))
}
