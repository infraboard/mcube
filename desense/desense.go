package desense

import (
	"fmt"
	"reflect"
	"strings"
)

// musk:defualt,3,2
func MaskStruct(s any) error {
	if reflect.TypeOf(s).Kind() != reflect.Ptr {
		return fmt.Errorf("object must be an pointer")
	}

	v := reflect.ValueOf(s).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		switch fieldValue.Kind() {
		case reflect.Slice:
			for i := 0; i < fieldValue.Len(); i++ {
				MaskStruct(fieldValue.Index(i).Interface())
			}
		case reflect.Struct:
			MaskStruct(fieldValue.Addr().Interface())
		default:
			tag := t.Field(i).Tag.Get("mask")
			if tag == "" {
				continue
			}
			name, args := ParseStructTag(tag)
			desenfor := Get(name)
			if desenfor == nil {
				return fmt.Errorf("desenfor %s not found", name)
			}
			if vStr, ok := fieldValue.Interface().(string); ok {
				fieldValue.SetString(desenfor.DeSense(vStr, args...))
			}
		}
	}
	return nil
}

func ParseStructTag(v string) (name string, args []string) {
	muskLine := strings.Split(v, ",")
	name = "default"
	if len(muskLine) > 0 {
		dn := muskLine[0]
		if dn != "" {
			name = dn
		}
	}
	if len(muskLine) > 1 {
		args = muskLine[1:]
	}
	return
}
