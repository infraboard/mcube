package ioc

import (
	"fmt"
	"reflect"
)

// Get 泛型方式获取对象，提供编译期类型安全
// 使用示例:
//
//	db, err := ioc.Get[*datasource.DataSource](ioc.Config())
//	if err != nil {
//	    return err
//	}
func Get[T Object](store StoreUser, opts ...GetOption) (T, error) {
	var zero T

	// 获取类型信息
	t := reflect.TypeOf(zero)
	if t == nil {
		return zero, fmt.Errorf("cannot get type information for nil")
	}

	// 获取类型名称
	name := t.String()
	if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
		name = t.String()
	}

	// 从store获取对象
	obj := store.Get(name, opts...)
	if obj == nil {
		return zero, fmt.Errorf("object %s not found in store", name)
	}

	// 类型断言
	result, ok := obj.(T)
	if !ok {
		return zero, fmt.Errorf("type mismatch: want %T, got %T", zero, obj)
	}

	return result, nil
}

// MustGet 泛型方式获取对象，如果失败则panic
// 使用示例:
//
//	db := ioc.MustGet[*datasource.DataSource](ioc.Config())
//
// 注意: 仅在确定对象存在时使用，否则会panic
func MustGet[T Object](store StoreUser, opts ...GetOption) T {
	obj, err := Get[T](store, opts...)
	if err != nil {
		panic(fmt.Sprintf("MustGet failed: %v", err))
	}
	return obj
}
