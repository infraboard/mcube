package diff

import "reflect"

// NormalizeNilContainers 返回 v 的深拷贝，并将 struct 内 nil map/slice 归一化为空容器，避免 diff 误报。
func NormalizeNilContainers[T any](v T) T {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		var zero T
		return zero
	}
	out := reflect.New(rv.Type()).Elem()
	normalizeInto(out, rv)
	return out.Interface().(T)
}

// CompareNormalized 在对比前对 before/after 做 nil 容器归一化。
func CompareNormalized[T any](before, after T, opts ...*Options) []DiffRecord {
	return Compare(
		NormalizeNilContainers(before),
		NormalizeNilContainers(after),
		opts...,
	)
}

func normalizeInto(dst, src reflect.Value) {
	if !src.IsValid() {
		return
	}

	switch src.Kind() {
	case reflect.Ptr:
		if src.IsNil() {
			return
		}
		if dst.IsNil() {
			dst.Set(reflect.New(src.Type().Elem()))
		}
		normalizeInto(dst.Elem(), src.Elem())
		return
	case reflect.Interface:
		if src.IsNil() {
			return
		}
		elem := src.Elem()
		cp := reflect.New(elem.Type()).Elem()
		normalizeInto(cp, elem)
		dst.Set(cp)
		return
	case reflect.Struct:
		for i := 0; i < src.NumField(); i++ {
			field := src.Type().Field(i)
			if !field.IsExported() {
				continue
			}
			normalizeInto(dst.Field(i), src.Field(i))
		}
		return
	case reflect.Map:
		if src.IsNil() {
			dst.Set(reflect.MakeMap(src.Type()))
			return
		}
		dst.Set(reflect.MakeMapWithSize(src.Type(), src.Len()))
		for _, key := range src.MapKeys() {
			v := src.MapIndex(key)
			nv := reflect.New(v.Type()).Elem()
			normalizeInto(nv, v)
			dst.SetMapIndex(key, nv)
		}
		return
	case reflect.Slice:
		if src.IsNil() {
			dst.Set(reflect.MakeSlice(src.Type(), 0, 0))
			return
		}
		dst.Set(reflect.MakeSlice(src.Type(), src.Len(), src.Len()))
		for i := 0; i < src.Len(); i++ {
			normalizeInto(dst.Index(i), src.Index(i))
		}
		return
	case reflect.Array:
		for i := 0; i < src.Len(); i++ {
			normalizeInto(dst.Index(i), src.Index(i))
		}
		return
	default:
		if dst.CanSet() {
			dst.Set(src)
		}
	}
}
