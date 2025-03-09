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
		case reflect.Struct:
			// 修正循环：使用正确的for循环遍历字段
			for i := range elems.NumField() {
				field := elems.Type().Field(i)
				if field.IsExported() {
					err := DensenceFiled(elems.Field(i), field.Tag)
					if err != nil {
						return err
					}
				}
			}
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
	for i := range elems.Len() {
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
		for i := range fieldValue.Len() {
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
	case reflect.Struct:
		// 处理嵌套的结构体值类型
		if fieldValue.CanAddr() {
			// 可寻址时，直接获取指针处理
			err := MaskStruct(fieldValue.Addr().Interface())
			if err != nil {
				return err
			}
		} else {
			// 不可寻址时，创建副本处理并回写
			copy := reflect.New(fieldValue.Type()).Elem()
			copy.Set(fieldValue)
			err := MaskStruct(copy.Addr().Interface())
			if err != nil {
				return err
			}
			fieldValue.Set(copy)
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
