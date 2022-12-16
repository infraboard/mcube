package exception

import "encoding/json"

// APIException API异常
type APIException interface {
	error
	ErrorCode() int
	WithHttpCode(int)
	HttpCode() int
	WithMeta(m interface{}) APIException
	Meta() interface{}
	WithData(d interface{}) APIException
	Data() interface{}
	Is(error) bool
	Namespace() string
	Reason() string
}

// APIException is impliment for api exception
type exception struct {
	namespace string
	httpCode  int
	code      int
	reason    string
	message   string
	meta      interface{}
	data      interface{}
}

func (e *exception) ToJson() string {
	m := map[string]interface{}{
		"namespace":  e.namespace,
		"http_code":  e.httpCode,
		"error_code": e.code,
		"reason":     e.reason,
		"meta":       e.meta,
		"data":       e.data,
		"message":    e.message,
	}
	dj, _ := json.Marshal(m)
	return string(dj)
}

func (e *exception) Error() string {
	return e.ToJson()
}

// Code exception's code, 如果code不存在返回-1
func (e *exception) ErrorCode() int {
	return int(e.code)
}

func (e *exception) WithHttpCode(httpCode int) {
	e.httpCode = httpCode
}

// Code exception's code, 如果code不存在返回-1
func (e *exception) HttpCode() int {
	return int(e.httpCode)
}

// WithMeta 携带一些额外信息
func (e *exception) WithMeta(m interface{}) APIException {
	e.meta = m
	return e
}

func (e *exception) Meta() interface{} {
	return e.meta
}

func (e *exception) WithData(d interface{}) APIException {
	e.data = d
	return e
}

func (e *exception) Data() interface{} {
	return e.data
}

func (e *exception) Is(t error) bool {
	if v, ok := t.(APIException); ok {
		return e.ErrorCode() == v.ErrorCode()
	}

	return e.message == t.Error()
}

func (e *exception) Namespace() string {
	return e.namespace
}

func (e *exception) Reason() string {
	return e.reason
}
