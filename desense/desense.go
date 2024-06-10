package desense

import (
	"fmt"
	"reflect"
	"strings"
)

// musk:defualt,3,2
func MaskStruct(target any) error {
	v := reflect.ValueOf(target)
	switch v.Kind() {
	case reflect.Ptr:
		elems := v.Elem()
		switch elems.Kind() {
		// 结构体指针 {}
		case reflect.Struct:
			for i := 0; i < elems.NumField(); i++ {
				err := DensenceFiled(elems.Field(i), elems.Type().Field(i).Tag)
				if err != nil {
					return err
				}
			}
		// 切片或者数组对象 指针
		case reflect.Slice, reflect.Array:
			err := DensenceList(elems)
			if err != nil {
				return err
			}
		}
	// 切片或者数组对象
	case reflect.Slice, reflect.Array:
		err := DensenceList(v)
		if err != nil {
			return err
		}
	}
	return nil
}

// DensenceFiled
func DensenceList(elems reflect.Value) error {
	for i := 0; i < elems.Len(); i++ {
		err := MaskStruct(elems.Index(i).Interface())
		if err != nil {
			return err
		}
	}

	return nil
}

func DensenceFiled(fieldValue reflect.Value, filedTag reflect.StructTag) error {
	switch fieldValue.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < fieldValue.Len(); i++ {
			err := MaskStruct(fieldValue.Index(i).Interface())
			if err != nil {
				return err
			}
		}
	case reflect.Ptr:
		err := MaskStruct(fieldValue.Interface())
		if err != nil {
			return err
		}
	case reflect.String:
		tag := filedTag.Get("mask")
		if tag == "" {
			return nil
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
