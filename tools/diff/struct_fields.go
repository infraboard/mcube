package diff

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type taggedScalarField struct {
	path []int
	name string
	desc string
	info fieldDiffTag
}

type sliceDisplayField struct {
	path []int
	Name string
	Desc string
}

func derefType(t reflect.Type) reflect.Type {
	for t != nil && t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func derefValue(v reflect.Value) reflect.Value {
	for v.IsValid() && v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return v
		}
		v = v.Elem()
	}
	return v
}

func fieldByPath(v reflect.Value, path []int) reflect.Value {
	for _, i := range path {
		v = derefValue(v)
		if !v.IsValid() || v.Kind() != reflect.Struct {
			return reflect.Value{}
		}
		v = v.Field(i)
	}
	return v
}

func structHasKeyField(t reflect.Type) bool {
	t = derefType(t)
	if t == nil || t.Kind() != reflect.Struct {
		return false
	}
	var walk func(rt reflect.Type) bool
	walk = func(rt reflect.Type) bool {
		rt = derefType(rt)
		if rt == nil || rt.Kind() != reflect.Struct {
			return false
		}
		for i := 0; i < rt.NumField(); i++ {
			field := rt.Field(i)
			if field.Tag.Get("diff") == "-" {
				continue
			}
			if field.Anonymous {
				if walk(field.Type) {
					return true
				}
				continue
			}
			if parseFieldDiffTag(field).isKey {
				return true
			}
		}
		return false
	}
	return walk(t)
}

func collectTaggedScalarFields(val reflect.Value) []taggedScalarField {
	if !val.IsValid() {
		return nil
	}
	val = derefValue(val)
	if !val.IsValid() || val.Kind() != reflect.Struct {
		return nil
	}
	var out []taggedScalarField
	var walk func(v reflect.Value, path []int, namePrefix, descPrefix string)
	walk = func(v reflect.Value, path []int, namePrefix, descPrefix string) {
		v = derefValue(v)
		if !v.IsValid() || v.Kind() != reflect.Struct {
			return
		}
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tag := field.Tag.Get("diff")
			if tag == "-" {
				continue
			}
			curPath := append(append([]int(nil), path...), i)
			fv := v.Field(i)
			if field.Anonymous {
				prefixName := namePrefix
				prefixDesc := descPrefix
				info := parseFieldDiffTag(field)
				if info.desc != "" && info.desc != field.Name {
					if prefixDesc != "" {
						prefixDesc += "."
					}
					prefixDesc += info.desc
				}
				walk(fv, curPath, prefixName, prefixDesc)
				continue
			}
			ft := field.Type
			for ft.Kind() == reflect.Ptr {
				ft = ft.Elem()
			}
			if !isScalarKind(ft.Kind()) {
				continue
			}
			if tag == "" {
				continue
			}
			info := parseFieldDiffTag(field)
			name := field.Name
			if namePrefix != "" {
				name = namePrefix + "." + name
			}
			desc := info.desc
			if descPrefix != "" {
				desc = descPrefix + "." + desc
			}
			out = append(out, taggedScalarField{path: curPath, name: name, desc: desc, info: info})
		}
	}
	walk(val, nil, "", "")
	return out
}

func getSliceElemDisplayFields(t reflect.Type) []sliceDisplayField {
	t = derefType(t)
	if t == nil || t.Kind() != reflect.Struct {
		return nil
	}
	useKeyOnly := structHasKeyField(t)
	zero := reflect.Zero(t)
	scalars := collectTaggedScalarFields(zero)
	fields := make([]sliceDisplayField, 0, len(scalars))
	for _, sf := range scalars {
		if useKeyOnly && !sf.info.isKey {
			continue
		}
		fields = append(fields, sliceDisplayField{path: sf.path, Name: sf.name, Desc: sf.desc})
	}
	return fields
}

func collectFieldValues(elems []reflect.Value, path []int) []string {
	vals := make([]string, 0, len(elems))
	for _, elem := range elems {
		if !elem.IsValid() {
			continue
		}
		fv := derefValue(fieldByPath(elem, path))
		if !fv.IsValid() {
			continue
		}
		vals = append(vals, fmt.Sprintf("%v", fv.Interface()))
	}
	return vals
}

func formatElemSnapshot(v reflect.Value) string {
	if !v.IsValid() {
		return ""
	}
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "<nil>"
		}
		v = v.Elem()
	}
	if b, err := json.Marshal(v.Interface()); err == nil {
		return string(b)
	}
	return fmt.Sprintf("%v", v.Interface())
}
