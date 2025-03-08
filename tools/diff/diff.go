package diff

import (
	"fmt"
	"reflect"
	"strings"
	"time"
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

type Options struct {
	IgnoreFields map[string]bool      // 忽略字段路径
	FieldLevels  map[string]DiffLevel // 字段级别映射
}

// 增强版对比入口
func Compare(a, b any, opts ...*Options) []DiffRecord {
	opt := &Options{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)

	if aVal.Type() != bVal.Type() {
		panic("comparison type mismatch")
	}

	return compare(aVal, bVal, "", "", opt, make([]DiffRecord, 0))
}

// 递归对比核心
func compare(aVal, bVal reflect.Value, path, descPath string, opt *Options, diffs []DiffRecord) []DiffRecord {
	// 处理指针
	if aVal.Kind() == reflect.Ptr {
		return comparePointer(aVal, bVal, path, descPath, opt, diffs)
	}

	// 基本类型直接比较（包含 time.Time）
	if isBasicType(aVal) {
		return compareBasic(aVal, bVal, path, descPath, opt, diffs)
	}

	// 根据类型分发处理
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

// 修改后的指针处理逻辑
func comparePointer(aVal, bVal reflect.Value, path, descPath string, opt *Options, diffs []DiffRecord) []DiffRecord {
	// 记录当前指针层级的路径
	ptrPath := path
	ptrDesc := descPath

	// 处理nil情况
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

	// 递归比较解引用后的值（保留指针层级信息）
	return compare(
		aVal.Elem(),
		bVal.Elem(),
		ptrPath,
		ptrDesc,
		opt,
		diffs,
	)
}

// 改进的标签解析
func parseDiffTag(field reflect.StructField) (desc string, level DiffLevel) {
	tag := field.Tag.Get("diff")
	if tag == "-" {
		return "", LevelInfo // 忽略字段时返回空描述
	}

	if tag == "" {
		return field.Name, LevelInfo
	}

	// 解析描述和级别
	parts := strings.Split(tag, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "level=") {
			levelStr := strings.TrimPrefix(part, "level=")
			switch strings.ToLower(levelStr) {
			case "warn":
				level = LevelWarn
			case "error":
				level = LevelError
			default:
				level = LevelInfo
			}
		} else if part != "" {
			desc = part // 第一个非 level 部分作为描述
		}
	}
	return desc, level
}

// 基本类型比较（包含 time.Time）
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
		fieldPath := buildPath(path, field)
		fieldDesc, level := buildDesc(descPath, field)

		// 补充字段级别
		if opt == nil {
			opt = &Options{
				FieldLevels: make(map[string]DiffLevel),
			}
		}
		if opt.FieldLevels == nil {
			opt.FieldLevels = make(map[string]DiffLevel)
		}
		opt.FieldLevels[fieldPath] = level

		// 检查是否忽略该字段（通过标签或配置）
		if shouldIgnoreField(field, fieldPath, opt) {
			continue
		}

		// 递归比较子字段
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

// 检查是否需要忽略字段
func shouldIgnoreField(field reflect.StructField, fieldPath string, opt *Options) bool {
	// 检查标签是否为 "-"
	tag := field.Tag.Get("diff")
	if tag == "-" {
		return true
	}

	// 检查是否在 IgnoreFields 配置中
	if opt != nil && opt.IgnoreFields != nil && opt.IgnoreFields[fieldPath] {
		return true
	}

	return false
}

// 切片/数组处理
func compareSlice(aVal, bVal reflect.Value, path, descPath string, opt *Options, diffs []DiffRecord) []DiffRecord {
	maxLen := max(aVal.Len(), bVal.Len())

	for i := range maxLen {
		elemPath := fmt.Sprintf("%s[%d]", path, i)

		// 安全获取元素值
		var aElem, bElem reflect.Value
		validA, validB := false, false

		if i < aVal.Len() {
			aElem = aVal.Index(i)
			validA = true
		}
		if i < bVal.Len() {
			bElem = bVal.Index(i)
			validB = true
		}

		// 处理越界情况
		if !validA && !validB {
			continue
		}

		// 如果元素是结构体或指针，递归比较
		if validA && validB {
			if aElem.Kind() == reflect.Struct || (aElem.Kind() == reflect.Ptr && aElem.Type().Elem().Kind() == reflect.Struct) {
				diffs = compare(aElem, bElem, elemPath, descPath, opt, diffs)
				continue
			}
		}

		// 处理新增或删除的情况
		if validA != validB {
			var oldVal, newVal any
			if validA {
				oldVal = aElem.Interface()
			}
			if validB {
				newVal = bElem.Interface()
			}

			// 如果是结构体，逐字段记录差异
			if (validA && aElem.Kind() == reflect.Struct) || (validB && bElem.Kind() == reflect.Struct) {
				diffs = compareStructElement(aElem, bElem, elemPath, descPath, opt, diffs)
			} else {
				// 非结构体类型，直接记录差异
				diffs = append(diffs, createDiff(
					elemPath,
					descPath,
					oldVal,
					newVal,
					opt,
				))
			}
			continue
		}

		// 转换为可安全调用的值
		var oldVal, newVal any
		if validA {
			oldVal = aElem.Interface()
		}
		if validB {
			newVal = bElem.Interface()
		}

		// 手动比较避免递归调用
		if !reflect.DeepEqual(oldVal, newVal) {
			diffs = append(diffs, createDiff(
				elemPath,
				descPath,
				oldVal,
				newVal,
				opt,
			))
		}
	}
	return diffs
}

func compareStructElement(aElem, bElem reflect.Value, path, descPath string, opt *Options, diffs []DiffRecord) []DiffRecord {
	// 处理新增结构体
	if !aElem.IsValid() || aElem.IsZero() {
		t := bElem.Type()
		for i := range bElem.NumField() {
			field := t.Field(i)
			fieldPath := buildPath(path, field)
			fieldDesc, level := buildDesc(descPath, field)

			// 补充字段级别
			if opt == nil {
				opt = &Options{
					FieldLevels: make(map[string]DiffLevel),
				}
			}
			if opt.FieldLevels == nil {
				opt.FieldLevels = make(map[string]DiffLevel)
			}
			opt.FieldLevels[fieldPath] = level

			// 检查是否忽略该字段
			if shouldIgnoreField(field, fieldPath, opt) {
				continue
			}

			// 记录新增字段
			diffs = append(diffs, createDiff(
				fieldPath,
				fieldDesc,
				nil,
				bElem.Field(i).Interface(),
				opt,
			))
		}
		return diffs
	}

	// 处理删除结构体
	if !bElem.IsValid() || bElem.IsZero() {
		t := aElem.Type()
		for i := range aElem.NumField() {
			field := t.Field(i)
			fieldPath := buildPath(path, field)
			fieldDesc, level := buildDesc(descPath, field)

			// 补充字段级别
			if opt == nil {
				opt = &Options{
					FieldLevels: make(map[string]DiffLevel),
				}
			}
			if opt.FieldLevels == nil {
				opt.FieldLevels = make(map[string]DiffLevel)
			}
			opt.FieldLevels[fieldPath] = level

			// 检查是否忽略该字段
			if shouldIgnoreField(field, fieldPath, opt) {
				continue
			}

			// 记录删除字段
			diffs = append(diffs, createDiff(
				fieldPath,
				fieldDesc,
				aElem.Field(i).Interface(),
				nil,
				opt,
			))
		}
		return diffs
	}

	return diffs
}

// Map 处理
func compareMap(aVal, bVal reflect.Value, path, descPath string, opt *Options, diffs []DiffRecord) []DiffRecord {
	allKeys := make(map[any]struct{})
	for _, key := range aVal.MapKeys() {
		allKeys[key.Interface()] = struct{}{}
	}
	for _, key := range bVal.MapKeys() {
		allKeys[key.Interface()] = struct{}{}
	}

	for key := range allKeys {
		keyStr := fmt.Sprintf("%v", key)
		elemPath := fmt.Sprintf("%s[%s]", path, keyStr)

		// 安全获取值
		var aItem, bItem reflect.Value
		validA, validB := false, false

		aItem = aVal.MapIndex(reflect.ValueOf(key))
		if aItem.IsValid() && !aItem.IsZero() {
			validA = true
		}
		bItem = bVal.MapIndex(reflect.ValueOf(key))
		if bItem.IsValid() && !bItem.IsZero() {
			validB = true
		}

		// 处理不存在的值
		var oldVal, newVal any
		if validA {
			oldVal = aItem.Interface()
		}
		if validB {
			newVal = bItem.Interface()
		}

		// 直接比较值
		if !reflect.DeepEqual(oldVal, newVal) {
			diffs = append(diffs, createDiff(
				elemPath,
				descPath,
				oldVal,
				newVal,
				opt,
			))
		}
	}
	return diffs
}

// 工具函数
func unpackPointer(v reflect.Value) any {
	if v.IsNil() {
		return nil
	}
	return v.Elem().Interface()
}

// 增强的路径构建逻辑
func buildPath(base string, field reflect.StructField) string {
	if shouldIgnoreField(field, "", nil) {
		return base // 忽略字段时不添加到路径
	}
	if base == "" {
		return field.Name
	}

	// 处理指针层级的*标记
	if strings.HasSuffix(base, "*") {
		return base + field.Name
	}

	return base + "." + field.Name
}

func buildDesc(base string, field reflect.StructField) (desc string, level DiffLevel) {
	if shouldIgnoreField(field, "", nil) {
		return base, LevelInfo // 忽略字段时不添加到描述
	}
	tag, level := parseDiffTag(field)
	if base == "" {
		return tag, level
	}
	return base + "." + tag, level
}

func createDiff(path, desc string, oldVal, newVal any, opt *Options) DiffRecord {
	// 自动确定级别
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
