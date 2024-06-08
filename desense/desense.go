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
	typ := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := typ.Field(i).Tag.Get("mask")

		muskLine := strings.Split(tag, ",")
		desensorName := "default"
		if len(muskLine) > 0 {
			dn := muskLine[0]
			if dn != "" {
				desensorName = dn
			}
		}
		desensorParam := []string{}
		if len(muskLine) > 1 {
			desensorParam = muskLine[1:]
		}

		desenfor := Get(desensorName)
		if desenfor == nil {
			return fmt.Errorf("desenfor %s not found", desensorName)
		}

		if field.Kind() == reflect.Struct {
			MaskStruct(field.Addr().Interface())
		} else {
			fieldValue := v.Field(i).Interface()
			if vStr, ok := fieldValue.(string); ok {
				field.SetString(desenfor.DeSense(vStr, desensorParam...))
			}
		}

	}

	return nil
}
