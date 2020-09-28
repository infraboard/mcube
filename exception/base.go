package exception

import "fmt"

// APIException API异常
type APIException interface {
	error
	ErrorCode() int
	WithMeta(m interface{})
	Meta() interface{}
	WithData(d interface{})
	Data() interface{}
	Is(code int) bool
	Namespace() string
	Reason() string
}

func newException(namespace Namespace, code int, format string, a ...interface{}) *exception {
	return &exception{
		namespace: namespace,
		code:      code,
		reason:    codeReason(code),
		message:   fmt.Sprintf(format, a...),
	}
}

// APIException is impliment for api exception
type exception struct {
	namespace Namespace
	code      int
	reason    string
	message   string
	meta      interface{}
	data      interface{}
}

func (e *exception) Error() string {
	return e.message
}

// Code exception's code, 如果code不存在返回-1
func (e *exception) ErrorCode() int {
	return int(e.code)
}

// WithMeta 携带一些额外信息
func (e *exception) WithMeta(m interface{}) {
	e.meta = m
}

func (e *exception) Meta() interface{} {
	return e.meta
}

func (e *exception) WithData(d interface{}) {
	e.data = d
}

func (e *exception) Data() interface{} {
	return e.data
}

func (e *exception) Is(code int) bool {
	return e.ErrorCode() == code
}

func (e *exception) Namespace() string {
	return e.namespace.String()
}

func (e *exception) Reason() string {
	return e.reason
}
