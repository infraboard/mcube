// Code generated by github.com/infraboard/mcube/v2
// DO NOT EDIT

package example

import (
	"bytes"
	"fmt"
	"strings"
)

// ParseFOOFromString Parse FOO from string
func ParseFOOFromString(str string) (FOO, error) {
	key := strings.Trim(string(str), `"`)
	v, ok := FOO_value[strings.ToUpper(key)]
	if !ok {
		return 0, fmt.Errorf("unknown FOO: %s", str)
	}

	return FOO(v), nil
}

// Equal type compare
func (t FOO) Equal(target FOO) bool {
	return t == target
}

// IsIn todo
func (t FOO) IsIn(targets ...FOO) bool {
	for _, target := range targets {
		if t.Equal(target) {
			return true
		}
	}

	return false
}

// MarshalJSON todo
func (t FOO) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(strings.ToUpper(t.String()))
	b.WriteString(`"`)
	return b.Bytes(), nil
}

// UnmarshalJSON todo
func (t *FOO) UnmarshalJSON(b []byte) error {
	ins, err := ParseFOOFromString(string(b))
	if err != nil {
		return err
	}
	*t = ins
	return nil
}
