// Package diff 提供基于 struct diff 标签的结构体字段对比，用于审计日志与变更事件。
//
// # 约定
//
//   - 参与对比的字段应显式标注 diff:"描述"；不参与对比用 diff:"-"；未标注字段默认也会对比（可用 Options.OnlyTaggedFields 关闭）。
//   - Compare / CompareNormalized 要求 before 与 after 为同一类型；CompareE / CompareNormalizedE 在类型不一致时返回 ErrTypeMismatch。
//   - 切片/数组对 struct（含 *struct）元素：先用 identity 字段匹配同一条元素，再递归比较其余字段。
//     标量字段可加 diff:"ID,key" 仅作 identity（同 ID 改其它字段会记 in-place diff，而非整行新增/移除）；
//     未标注 key 时退化为「全部带 diff 标签的标量字段拼接」作为 identity。
//     匹配后的字段路径会带 identity/下标，如 Items[1].Name。
//   - Map 按 key 对齐后对 value 递归对比，路径形如 Meta[env].Name。
//
// # 典型用法
//
//	records := diff.CompareNormalized(beforeSpec, afterSpec)
//	records = diff.AppendRecords(records, diff.MaskedChange("Value", "凭证值", "******", diff.LevelWarn))
package diff

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	// MarkerAdd 切片/集合新增元素的占位旧值。
	MarkerAdd = "新增"
	// MarkerRemove 切片/集合移除元素的占位新值。
	MarkerRemove = "移除"
)

type DiffLevel int

const (
	LevelInfo DiffLevel = iota
	LevelWarn
	LevelError
)

type DiffRecord struct {
	FieldPath string // 结构体字段路径（如 "Basic.ID"）
	FieldDesc string // 字段描述（如 "基础信息.用户ID"）
	OldValue  any
	NewValue  any
	Level     DiffLevel
}

func (s *DiffRecord) GetOldValue() string {
	return formatValue(s.OldValue)
}

func (s *DiffRecord) GetNewValue() string {
	return formatValue(s.NewValue)
}

func formatValue(v any) string {
	if v == nil {
		return "N/A"
	}
	if str, ok := v.(string); ok {
		return str
	}
	switch val := v.(type) {
	case time.Time:
		return val.Format(time.RFC3339)
	case []byte:
		return string(val)
	default:
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Map, reflect.Slice, reflect.Array, reflect.Struct:
			if b, err := json.Marshal(v); err == nil {
				return string(b)
			}
		}
	}
	return fmt.Sprintf("%v", v)
}

type Options struct {
	IgnoreFields map[string]bool      // 忽略字段路径
	FieldLevels  map[string]DiffLevel // 字段级别映射
	// OnlyTaggedFields 为 true 时，仅对比显式带 diff 标签（且非 "-"）的字段。
	OnlyTaggedFields bool
}

// Compare 对比两个同类型对象，返回差异列表；类型不一致时 panic。
func Compare(a, b any, opts ...*Options) []DiffRecord {
	records, err := CompareE(a, b, opts...)
	if err != nil {
		panic(err)
	}
	return records
}

// Summary 将多条差异格式化为可读摘要（field: old -> new）。
func Summary(records []DiffRecord) string {
	if len(records) == 0 {
		return ""
	}
	parts := make([]string, 0, len(records))
	for _, r := range records {
		parts = append(parts, fmt.Sprintf("%s: %s -> %s", r.FieldDesc, r.GetOldValue(), r.GetNewValue()))
	}
	return strings.Join(parts, ", ")
}

func compare(aVal, bVal reflect.Value, path, descPath string, opt *Options, diffs []DiffRecord) []DiffRecord {
	if !aVal.IsValid() || !bVal.IsValid() {
		if aVal.IsValid() != bVal.IsValid() {
			var oldVal, newVal any
			if aVal.IsValid() {
				oldVal = aVal.Interface()
			}
			if bVal.IsValid() {
				newVal = bVal.Interface()
			}
			return append(diffs, createDiff(path, descPath, oldVal, newVal, opt))
		}
		return diffs
	}

	switch aVal.Kind() {
	case reflect.Interface:
		return compareInterface(aVal, bVal, path, descPath, opt, diffs)
	case reflect.Ptr:
		return comparePointer(aVal, bVal, path, descPath, opt, diffs)
	}

	if isBasicType(aVal) {
		return compareBasic(aVal, bVal, path, descPath, opt, diffs)
	}

	switch aVal.Kind() {
	case reflect.Struct:
		return compareStruct(aVal, bVal, path, descPath, opt, diffs)
	case reflect.Slice, reflect.Array:
		return compareSlice(aVal, bVal, path, descPath, opt, diffs)
	case reflect.Map:
		return compareMap(aVal, bVal, path, descPath, opt, diffs)
	default:
		return diffs
	}
}

func compareInterface(aVal, bVal reflect.Value, path, descPath string, opt *Options, diffs []DiffRecord) []DiffRecord {
	if aVal.IsNil() && bVal.IsNil() {
		return diffs
	}
	if aVal.IsNil() || bVal.IsNil() {
		var oldVal, newVal any
		if !aVal.IsNil() {
			oldVal = aVal.Elem().Interface()
		}
		if !bVal.IsNil() {
			newVal = bVal.Elem().Interface()
		}
		return append(diffs, createDiff(path, descPath, oldVal, newVal, opt))
	}

	aElem, bElem := aVal.Elem(), bVal.Elem()
	if aElem.Type() != bElem.Type() {
		return append(diffs, createDiff(path, descPath, aElem.Interface(), bElem.Interface(), opt))
	}
	return compare(aElem, bElem, path, descPath, opt, diffs)
}

func comparePointer(aVal, bVal reflect.Value, path, descPath string, opt *Options, diffs []DiffRecord) []DiffRecord {
	ptrPath := path
	ptrDesc := descPath

	if aVal.IsNil() || bVal.IsNil() {
		if aVal.IsNil() != bVal.IsNil() {
			diffs = append(diffs, createDiff(
				path,
				descPath,
				unpackPointer(aVal),
				unpackPointer(bVal),
				opt,
			))
		}
		return diffs
	}

	return compare(
		aVal.Elem(),
		bVal.Elem(),
		ptrPath,
		ptrDesc,
		opt,
		diffs,
	)
}

func parseDiffTag(field reflect.StructField) (desc string, level DiffLevel) {
	info := parseFieldDiffTag(field)
	return info.desc, info.level
}

type fieldDiffTag struct {
	desc  string
	level DiffLevel
	isKey bool
}

func parseFieldDiffTag(field reflect.StructField) fieldDiffTag {
	tag := field.Tag.Get("diff")
	info := fieldDiffTag{desc: field.Name, level: LevelInfo}
	if tag == "-" {
		return fieldDiffTag{}
	}
	if tag == "" {
		return info
	}

	parts := strings.Split(tag, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		switch {
		case strings.HasPrefix(part, "level="):
			switch strings.ToLower(strings.TrimPrefix(part, "level=")) {
			case "warn":
				info.level = LevelWarn
			case "error":
				info.level = LevelError
			default:
				info.level = LevelInfo
			}
		case part == "key":
			info.isKey = true
		case part != "":
			info.desc = part
		}
	}
	return info
}

func isScalarKind(kind reflect.Kind) bool {
	return kind == reflect.String || kind == reflect.Bool ||
		(kind >= reflect.Int && kind <= reflect.Int64) ||
		(kind >= reflect.Uint && kind <= reflect.Uint64) ||
		(kind >= reflect.Float32 && kind <= reflect.Float64)
}

func isBasicType(v reflect.Value) bool {
	if v.Type() == reflect.TypeOf(time.Time{}) {
		return true
	}
	switch v.Kind() {
	case reflect.String, reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func compareBasic(aVal, bVal reflect.Value, path, descPath string, opt *Options, diffs []DiffRecord) []DiffRecord {
	if !reflect.DeepEqual(aVal.Interface(), bVal.Interface()) {
		return append(diffs, createDiff(path, descPath,
			aVal.Interface(),
			bVal.Interface(),
			opt))
	}
	return diffs
}

func compareStruct(aVal, bVal reflect.Value, path, descPath string, opt *Options, diffs []DiffRecord) []DiffRecord {
	t := aVal.Type()
	for i := range aVal.NumField() {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		fieldPath := buildPath(path, field)
		fieldDesc, level := buildDesc(descPath, field)

		if opt == nil {
			opt = &Options{FieldLevels: make(map[string]DiffLevel)}
		}
		if opt.FieldLevels == nil {
			opt.FieldLevels = make(map[string]DiffLevel)
		}
		// 仅在调用方未指定时写入 tag 默认级别，避免覆盖 Options.FieldLevels
		if _, exists := opt.FieldLevels[fieldPath]; !exists {
			opt.FieldLevels[fieldPath] = level
		}

		if shouldIgnoreField(field, fieldPath, opt) {
			continue
		}

		diffs = compare(
			aVal.Field(i),
			bVal.Field(i),
			fieldPath,
			fieldDesc,
			opt,
			diffs,
		)
	}
	return diffs
}

func shouldIgnoreField(field reflect.StructField, fieldPath string, opt *Options) bool {
	tag := field.Tag.Get("diff")
	if tag == "-" {
		return true
	}
	if opt != nil && opt.OnlyTaggedFields && tag == "" {
		return true
	}
	if opt != nil && opt.IgnoreFields != nil && opt.IgnoreFields[fieldPath] {
		return true
	}
	return false
}

func compareSlice(aVal, bVal reflect.Value, path, descPath string, opt *Options, diffs []DiffRecord) []DiffRecord {
	usedB := make([]bool, bVal.Len())
	type indexPair struct{ ai, bi int }
	matched := make([]indexPair, 0, minInt(aVal.Len(), bVal.Len()))

	for i := 0; i < aVal.Len(); i++ {
		aKey := sliceElemKey(aVal.Index(i))
		for j := 0; j < bVal.Len(); j++ {
			if usedB[j] {
				continue
			}
			if sliceElemKey(bVal.Index(j)) == aKey {
				usedB[j] = true
				matched = append(matched, indexPair{i, j})
				break
			}
		}
	}

	for _, pair := range matched {
		elemPath, elemDesc := sliceElemPath(path, descPath, aVal.Index(pair.ai), pair.ai)
		diffs = compare(aVal.Index(pair.ai), bVal.Index(pair.bi), elemPath, elemDesc, opt, diffs)
	}

	var removedElems []reflect.Value
	for i := 0; i < aVal.Len(); i++ {
		found := false
		for _, pair := range matched {
			if pair.ai == i {
				found = true
				break
			}
		}
		if !found {
			removedElems = append(removedElems, aVal.Index(i))
		}
	}

	var addedElems []reflect.Value
	for j := 0; j < bVal.Len(); j++ {
		if !usedB[j] {
			addedElems = append(addedElems, bVal.Index(j))
		}
	}

	if len(removedElems) == 0 && len(addedElems) == 0 {
		return diffs
	}

	diffs = appendSliceElemDiffs(aVal, bVal, path, descPath, opt, diffs, removedElems, addedElems)
	return diffs
}

func appendSliceElemDiffs(
	aVal, bVal reflect.Value,
	path, descPath string,
	opt *Options,
	diffs []DiffRecord,
	removedElems, addedElems []reflect.Value,
) []DiffRecord {
	var elemType reflect.Type
	if aVal.Len() > 0 {
		elemType = aVal.Index(0).Type()
	} else if bVal.Len() > 0 {
		elemType = bVal.Index(0).Type()
	}

	structType := derefType(elemType)
	if structType != nil && structType.Kind() == reflect.Struct {
		displayFields := getSliceElemDisplayFields(structType)
		if len(displayFields) > 0 {
			for _, df := range displayFields {
				subPath := path + "." + df.Name
				subDescPath := descPath + "." + df.Desc
				if len(removedElems) > 0 {
					for _, val := range collectFieldValues(removedElems, df.path) {
						diffs = append(diffs, createDiff(subPath, subDescPath, val, MarkerRemove, opt))
					}
				}
				if len(addedElems) > 0 {
					for _, val := range collectFieldValues(addedElems, df.path) {
						diffs = append(diffs, createDiff(subPath, subDescPath, MarkerAdd, val, opt))
					}
				}
			}
			return diffs
		}
	}

	for _, e := range removedElems {
		elemPath, elemDesc := sliceElemPath(path, descPath, e, 0)
		diffs = append(diffs, createDiff(elemPath, elemDesc, formatElemSnapshot(e), MarkerRemove, opt))
	}
	for _, e := range addedElems {
		elemPath, elemDesc := sliceElemPath(path, descPath, e, 0)
		diffs = append(diffs, createDiff(elemPath, elemDesc, MarkerAdd, formatElemSnapshot(e), opt))
	}
	return diffs
}

// sliceElemPath 为切片元素生成带 identity/下标的路径，避免多条元素变更路径冲突。
func sliceElemPath(path, descPath string, elem reflect.Value, index int) (string, string) {
	label := sliceElemPathLabel(elem, index)
	suffix := "[" + label + "]"
	if path == "" {
		return suffix, descPath + suffix
	}
	return path + suffix, descPath + suffix
}

func sliceElemPathLabel(elem reflect.Value, index int) string {
	v := derefValue(elem)
	if !v.IsValid() {
		return strconv.Itoa(index)
	}
	if v.Kind() == reflect.Struct {
		if structHasKeyField(v.Type()) {
			if key := sliceElemKey(elem); key != "" {
				return key
			}
		}
		return strconv.Itoa(index)
	}
	return strconv.Itoa(index)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func sliceElemKey(v reflect.Value) string {
	if !v.IsValid() {
		return ""
	}
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "<nil>"
		}
		return sliceElemKey(v.Elem())
	}
	if isBasicType(v) {
		return fmt.Sprintf("%v", v.Interface())
	}
	if v.Kind() == reflect.Struct {
		t := v.Type()
		scalars := collectTaggedScalarFields(v)
		useKeyOnly := structHasKeyField(t)
		var keyParts, fallbackParts []string
		for _, sf := range scalars {
			fv := fieldByPath(v, sf.path)
			fv = derefValue(fv)
			if !fv.IsValid() {
				continue
			}
			val := fmt.Sprintf("%v", fv.Interface())
			if useKeyOnly {
				if sf.info.isKey {
					keyParts = append(keyParts, val)
				}
				continue
			}
			fallbackParts = append(fallbackParts, val)
		}
		if useKeyOnly && len(keyParts) > 0 {
			return strings.Join(keyParts, "|")
		}
		if len(fallbackParts) > 0 {
			return strings.Join(fallbackParts, "|")
		}
	}
	return formatElemSnapshot(v)
}

func compareMap(aVal, bVal reflect.Value, path, descPath string, opt *Options, diffs []DiffRecord) []DiffRecord {
	type keyRef struct {
		rv reflect.Value
		id any
	}
	seen := make(map[any]struct{})
	keys := make([]keyRef, 0, aVal.Len()+bVal.Len())
	addKey := func(k reflect.Value) {
		id := k.Interface()
		if _, ok := seen[id]; ok {
			return
		}
		seen[id] = struct{}{}
		keys = append(keys, keyRef{rv: k, id: id})
	}
	for _, key := range aVal.MapKeys() {
		addKey(key)
	}
	for _, key := range bVal.MapKeys() {
		addKey(key)
	}

	for _, key := range keys {
		keyStr := fmt.Sprintf("%v", key.id)
		elemPath := path
		if elemPath == "" {
			elemPath = "[" + keyStr + "]"
		} else {
			elemPath = fmt.Sprintf("%s[%s]", path, keyStr)
		}
		elemDesc := descPath
		if descPath != "" {
			elemDesc = descPath + "." + keyStr
		} else {
			elemDesc = keyStr
		}

		aItem := aVal.MapIndex(key.rv)
		bItem := bVal.MapIndex(key.rv)

		switch {
		case aItem.IsValid() && bItem.IsValid():
			diffs = compare(aItem, bItem, elemPath, elemDesc, opt, diffs)
		case aItem.IsValid():
			diffs = append(diffs, createDiff(elemPath, elemDesc, aItem.Interface(), nil, opt))
		case bItem.IsValid():
			diffs = append(diffs, createDiff(elemPath, elemDesc, nil, bItem.Interface(), opt))
		}
	}
	return diffs
}

func unpackPointer(v reflect.Value) any {
	if v.IsNil() {
		return nil
	}
	return v.Elem().Interface()
}

func buildPath(base string, field reflect.StructField) string {
	if field.Tag.Get("diff") == "-" {
		return base
	}
	if base == "" {
		return field.Name
	}
	if strings.HasSuffix(base, "*") {
		return base + field.Name
	}
	return base + "." + field.Name
}

func buildDesc(base string, field reflect.StructField) (desc string, level DiffLevel) {
	if field.Tag.Get("diff") == "-" {
		return base, LevelInfo
	}
	tag, level := parseDiffTag(field)
	if base == "" {
		return tag, level
	}
	return base + "." + tag, level
}

func createDiff(path, desc string, oldVal, newVal any, opt *Options) DiffRecord {
	level := LevelInfo
	if opt != nil && opt.FieldLevels != nil {
		if lvl, ok := opt.FieldLevels[path]; ok {
			level = lvl
		}
	}

	return DiffRecord{
		FieldPath: path,
		FieldDesc: desc,
		OldValue:  oldVal,
		NewValue:  newVal,
		Level:     level,
	}
}
