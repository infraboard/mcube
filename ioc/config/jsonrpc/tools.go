package jsonrpc

import (
	"reflect"
	"runtime"
	"strings"
)

// 获取函数名
func getFunctionName(fn HandlerFunc) string {
	pc := reflect.ValueOf(fn).Pointer()
	fullName := runtime.FuncForPC(pc).Name()

	// 简化函数名，去掉包路径等冗余信息
	parts := strings.Split(fullName, ".")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return fullName
}

// 获取类型名
func getTypeName(obj any) string {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}
