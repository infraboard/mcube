// Code generated by github.com/infraboard/mcube
// DO NOT EDIT

package resource

import (
	"bytes"
	"fmt"
	"strings"
)

// ParseVISIABLEFromString Parse VISIABLE from string
func ParseVISIABLEFromString(str string) (VISIABLE, error) {
	key := strings.Trim(string(str), `"`)
	v, ok := VISIABLE_value[strings.ToUpper(key)]
	if !ok {
		return 0, fmt.Errorf("unknown VISIABLE: %s", str)
	}

	return VISIABLE(v), nil
}

// Equal type compare
func (t VISIABLE) Equal(target VISIABLE) bool {
	return t == target
}

// IsIn todo
func (t VISIABLE) IsIn(targets ...VISIABLE) bool {
	for _, target := range targets {
		if t.Equal(target) {
			return true
		}
	}

	return false
}

// MarshalJSON todo
func (t VISIABLE) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(strings.ToUpper(t.String()))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

// UnmarshalJSON todo
func (t *VISIABLE) UnmarshalJSON(b []byte) error {
	ins, err := ParseVISIABLEFromString(string(b))
	if err != nil {
		return err
	}
	*t = ins
	return nil
}