package desense_test

import (
	"testing"

	"github.com/infraboard/mcube/v2/desense"
)

func TestDeSense(t *testing.T) {
	v := desense.Default().DeSense("abcdefghijkl", "3", "2")
	t.Logf(v)
}

type TestMaskStuct struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number" mask:",5,4"`
}

func TestMusk(t *testing.T) {
	tm := &TestMaskStuct{
		PhoneNumber: "1234567890123",
		Name:        "test",
	}

	if err := desense.MaskStruct(tm); err != nil {
		t.Fatal(err)
	}
	t.Log(tm.PhoneNumber)
	t.Log(tm.Name)
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
	Items []*TestMaskStuct `json:"items"`
}

func TestMuskSet(t *testing.T) {
	if err := desense.MaskStruct(set); err != nil {
		t.Fatal(err)
	}

	for _, v := range set.Items {
		t.Log(v.PhoneNumber)
		t.Log(v.Name)
	}
}
