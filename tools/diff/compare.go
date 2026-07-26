package diff

import (
	"errors"
	"maps"
	"reflect"
)

// ErrTypeMismatch 在 CompareE / CompareNormalizedE 中，当 before 与 after 类型不一致时返回。
var ErrTypeMismatch = errors.New("diff: type mismatch")

// CompareE 与 Compare 相同，但类型不一致时返回 ErrTypeMismatch 而非 panic。
func CompareE(a, b any, opts ...*Options) ([]DiffRecord, error) {
	opt := &Options{}
	if len(opts) > 0 && opts[0] != nil {
		opt = cloneOptions(opts[0])
	}
	if opt.FieldLevels == nil {
		opt.FieldLevels = make(map[string]DiffLevel)
	}

	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)

	if !aVal.IsValid() || !bVal.IsValid() {
		if !aVal.IsValid() && !bVal.IsValid() {
			return nil, nil
		}
		return nil, ErrTypeMismatch
	}

	if aVal.Type() != bVal.Type() {
		return nil, ErrTypeMismatch
	}

	return compare(aVal, bVal, "", "", opt, make([]DiffRecord, 0)), nil
}

// CompareNormalizedE 归一化后对比，类型不一致时返回 ErrTypeMismatch。
func CompareNormalizedE[T any](before, after T, opts ...*Options) ([]DiffRecord, error) {
	return CompareE(
		NormalizeNilContainers(before),
		NormalizeNilContainers(after),
		opts...,
	)
}

func cloneOptions(opt *Options) *Options {
	out := *opt
	if opt.IgnoreFields != nil {
		out.IgnoreFields = make(map[string]bool, len(opt.IgnoreFields))
		maps.Copy(out.IgnoreFields, opt.IgnoreFields)
	}
	if opt.FieldLevels != nil {
		out.FieldLevels = make(map[string]DiffLevel, len(opt.FieldLevels))
		maps.Copy(out.FieldLevels, opt.FieldLevels)
	}
	return &out
}
