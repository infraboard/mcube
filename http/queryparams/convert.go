package queryparams

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// Marshaler converts an object to a query parameter string representation
type Marshaler interface {
	MarshalQueryParameter() (string, error)
}

// Unmarshaler converts a string representation to an object
type Unmarshaler interface {
	UnmarshalQueryParameter(string) error
}

func jsonTag(field reflect.StructField) (string, bool) {
	structTag := field.Tag.Get("json")
	if len(structTag) == 0 {
		return "", false
	}
	parts := strings.Split(structTag, ",")
	tag := parts[0]
	if tag == "-" {
		tag = ""
	}
	omitempty := false
	parts = parts[1:]
	for _, part := range parts {
		if part == "omitempty" {
			omitempty = true
			break
		}
	}
	return tag, omitempty
}

func isPointerKind(kind reflect.Kind) bool {
	return kind == reflect.Pointer
}

func isStructKind(kind reflect.Kind) bool {
	return kind == reflect.Struct
}

func isValueKind(kind reflect.Kind) bool {
	switch kind {
	case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8,
		reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32,
		reflect.Float64, reflect.Complex64, reflect.Complex128:
		return true
	default:
		return false
	}
}

func zeroValue(value reflect.Value) bool {
	return reflect.DeepEqual(reflect.Zero(value.Type()).Interface(), value.Interface())
}

func customMarshalValue(value reflect.Value) (reflect.Value, bool) {
	// Return unless we implement a custom query marshaler
	if !value.CanInterface() {
		return reflect.Value{}, false
	}

	marshaler, ok := value.Interface().(Marshaler)
	if !ok {
		if !isPointerKind(value.Kind()) && value.CanAddr() {
			marshaler, ok = value.Addr().Interface().(Marshaler)
			if !ok {
				return reflect.Value{}, false
			}
		} else {
			return reflect.Value{}, false
		}
	}

	// Don't invoke functions on nil pointers
	// If the type implements MarshalQueryParameter, AND the tag is not omitempty, AND the value is a nil pointer, "" seems like a reasonable response
	if isPointerKind(value.Kind()) && zeroValue(value) {
		return reflect.ValueOf(""), true
	}

	// Get the custom marshalled value
	v, err := marshaler.MarshalQueryParameter()
	if err != nil {
		return reflect.Value{}, false
	}
	return reflect.ValueOf(v), true
}

func addParam(values url.Values, tag string, omitempty bool, value reflect.Value) {
	if omitempty && zeroValue(value) {
		return
	}
	val := ""
	iValue := fmt.Sprintf("%v", value.Interface())

	if iValue != "<nil>" {
		val = iValue
	}
	values.Add(tag, val)
}

func addListOfParams(values url.Values, tag string, omitempty bool, list reflect.Value) {
	for i := 0; i < list.Len(); i++ {
		addParam(values, tag, omitempty, list.Index(i))
	}
}

// Convert takes an object and converts it to a url.Values object using JSON tags as
// parameter names. Only top-level simple values, arrays, and slices are serialized.
// not Embedded structs, maps, etc. will not be serialized.
func Convert(obj interface{}) (url.Values, error) {
	result := url.Values{}
	if obj == nil {
		return result, nil
	}

	var sv reflect.Value
	t := reflect.TypeOf(obj).Kind()
	switch t {
	case reflect.Pointer, reflect.Interface:
		sv = reflect.ValueOf(obj).Elem()
	case reflect.Struct:
		sv = reflect.ValueOf(obj)
	default:
		return nil, fmt.Errorf("expecting a pointer, struct or interface, but %s", t)
	}

	// Check all object fields
	convertStruct(result, sv)

	return result, nil
}

func convertStruct(result url.Values, sv reflect.Value) {
	st := sv.Type()

	for i := 0; i < st.NumField(); i++ {
		filed := st.Field(i)
		fieldValue := sv.Field(i)

		// 处理匿名嵌套
		if filed.Anonymous {
			switch fieldValue.Kind() {
			case reflect.Ptr:
				if fieldValue.Elem().Kind() == reflect.Struct {
					convertStruct(result, fieldValue.Elem())
				}
			case reflect.Struct:
				convertStruct(result, fieldValue)
			}
		}

		tag, omitempty := jsonTag(filed)
		if len(tag) == 0 {
			continue
		}
		ft := fieldValue.Type()
		kind := ft.Kind()
		if isPointerKind(kind) {
			ft = ft.Elem()
			kind = ft.Kind()
			if !fieldValue.IsNil() {
				fieldValue = reflect.Indirect(fieldValue)
				// If the field is non-nil, it should be added to params
				// and the omitempty should be overwite to false
				omitempty = false
			}
		}

		switch {
		case isValueKind(kind):
			addParam(result, tag, omitempty, fieldValue)
		case kind == reflect.Array || kind == reflect.Slice:
			if isValueKind(ft.Elem().Kind()) {
				addListOfParams(result, tag, omitempty, fieldValue)
			}
		case isStructKind(kind) && !(zeroValue(fieldValue) && omitempty):
			if marshalValue, ok := customMarshalValue(fieldValue); ok {
				addParam(result, tag, omitempty, marshalValue)
			} else {
				convertStruct(result, fieldValue)
			}
		}
	}
}
