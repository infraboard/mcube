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
	PhoneNumber string `json:"phone_number" mask:",5,4"`
}

func TestMusk(t *testing.T) {
	tm := &TestMaskStuct{
		PhoneNumber: "1234567890123",
	}
	if err := desense.MaskStruct(tm); err != nil {
		t.Fatal(err)
	}
	t.Log(tm.PhoneNumber)
}
