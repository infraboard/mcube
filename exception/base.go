package exception

import "fmt"

// APIException API异常
type APIException interface {
	error
	ErrorCode() int
	Meta() interface{}
	Is(code int) bool
}

// WithAPIException 携带元信息的异常
type WithAPIException interface {
	APIException
	WithMeta(m interface{}) APIException
}

// APIException is impliment for api exception
type exception struct {
	code int
	msg  string
	meta interface{}
}

func (e *exception) Error() string {
	return e.msg
}

// Code exception's code, 如果code不存在返回-1
func (e *exception) ErrorCode() int {
	return e.code
}

// WithMeta 携带一些额外信息
func (e *exception) WithMeta(m interface{}) APIException {
	e.meta = m
	return e
}

func (e *exception) Meta() interface{} {
	return e.meta
}

func (e *exception) Is(code int) bool {
	return e.ErrorCode() == code
}

// NewAPIException 创建一个API异常
// 用于其他模块自定义异常
func NewAPIException(code int, format string, a ...interface{}) WithAPIException {
	// 0表示正常状态, 但是要排除变量的零值
	if code == 0 {
		code = -1
	}

	return &exception{
		msg:  fmt.Sprintf(format, a...),
		code: code,
	}
}
