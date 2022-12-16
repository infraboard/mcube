package exception

import "encoding/json"

// APIException API异常
type APIException interface {
	error
	ErrorCode() int
	WithHttpCode(int)
	GetHttpCode() int
	WithMeta(m interface{}) APIException
	GetMeta() interface{}
	WithData(d interface{}) APIException
	GetData() interface{}
	Is(error) bool
	GetNamespace() string
	WithNamespace(ns string)
	GetReason() string
	ToJson() string
}

// APIException is impliment for api exception
type exception struct {
	Namespace string      `json:"namespace"`
	HttpCode  int         `json:"http_code"`
	Code      int         `json:"error_code"`
	Reason    string      `json:"reason"`
	Message   string      `json:"message"`
	Meta      interface{} `json:"meta"`
	Data      interface{} `json:"data"`
}

func (e *exception) ToJson() string {
	dj, _ := json.Marshal(e)
	return string(dj)
}

func (e *exception) Error() string {
	return e.Message
}

// Code exception's code, 如果code不存在返回-1
func (e *exception) ErrorCode() int {
	return int(e.Code)
}

func (e *exception) WithHttpCode(httpCode int) {
	e.HttpCode = httpCode
}

// Code exception's code, 如果code不存在返回-1
func (e *exception) GetHttpCode() int {
	return int(e.HttpCode)
}

// WithMeta 携带一些额外信息
func (e *exception) WithMeta(m interface{}) APIException {
	e.Meta = m
	return e
}

func (e *exception) GetMeta() interface{} {
	return e.Meta
}

func (e *exception) WithData(d interface{}) APIException {
	e.Data = d
	return e
}

func (e *exception) GetData() interface{} {
	return e.Data
}

func (e *exception) Is(t error) bool {
	if v, ok := t.(APIException); ok {
		return e.ErrorCode() == v.ErrorCode()
	}

	return e.Message == t.Error()
}

func (e *exception) GetNamespace() string {
	return e.Namespace
}

func (e *exception) GetReason() string {
	return e.Reason
}

func (e *exception) WithNamespace(ns string) {
	e.Namespace = ns
}
