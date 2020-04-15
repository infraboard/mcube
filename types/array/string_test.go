package array_test

import (
	"testing"

	"github.com/infraboard/mcube/types/array"
	"github.com/stretchr/testify/assert"
)

func TestStringArrayString(t *testing.T) {
	should := assert.New(t)
	data := array.NewStringArray([]string{"1", "2", "3"})
	should.Equal(data.String(), ";1;2;3;")
}

func TestStringArrayValue(t *testing.T) {
	should := assert.New(t)
	data := array.NewStringArray([]string{"1", "2", "3"})
	v, err := data.Value()
	if should.NoError(err) {
		strv, ok := v.(string)
		if should.True(ok) {
			should.Equal(strv, ";1;2;3;")
		}
	}
}

func TestStringArrayScan(t *testing.T) {
	should := assert.New(t)
	data := array.NewStringArray([]string{})
	if should.NoError(data.Scan([]byte(";1;2;3;"))) {
		should.Equal(data.Items(), []string{"1", "2", "3"})
	}
}
