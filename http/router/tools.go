package router

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

// GetHandlerFuncName 通过反射获取函数名称
func GetHandlerFuncName(h http.HandlerFunc, seps ...rune) string {
	// 获取函数名称
	fn := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()

	// 默认使用.分隔
	if len(seps) == 0 {
		seps = append(seps, '.')
	}

	// 用 seps 进行分割
	fields := strings.FieldsFunc(fn, func(sep rune) bool {
		for _, s := range seps {
			if sep == s {
				return true
			}
		}
		return false
	})

	if size := len(fields); size > 0 {
		return strings.TrimRight(fields[size-1], "-fm")
	}

	return ""
}
