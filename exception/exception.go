package exception

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	// GRPC Trailer 异常转换时定义的key名称
	TRAILER_ERROR_JSON_KEY = "err_json"
)

// NewAPIException 创建一个API异常
// 用于其他模块自定义异常
func NewAPIException(namespace string, code int, reason, format string, a ...interface{}) *APIException {
	// 0表示正常状态, 但是要排除变量的零值
	if code == 0 {
		code = -1
	}

	var httpCode int
	if code/100 >= 1 && code/100 <= 5 {
		httpCode = code
	} else {
		httpCode = http.StatusInternalServerError
	}

	return &APIException{
		Namespace: namespace,
		ErrCode:   code,
		Reason:    reason,
		HttpCode:  httpCode,
		Message:   fmt.Sprintf(format, a...),
	}
}

func NewAPIExceptionFromString(msg string) *APIException {
	e := &APIException{}
	if !strings.HasPrefix(msg, "{") {
		e.Message = msg
		e.ErrCode = InternalServerError
		e.HttpCode = InternalServerError
		return e
	}

	err := json.Unmarshal([]byte(msg), e)
	if err != nil {
		e.Message = msg
		e.ErrCode = InternalServerError
		e.HttpCode = InternalServerError
	}
	return e
}

// APIException API异常
type APIException struct {
	Namespace string `json:"namespace"`
	HttpCode  int    `json:"http_code"`
	ErrCode   int    `json:"error_code"`
	Reason    string `json:"reason"`
	Message   string `json:"message"`
	Meta      any    `json:"meta"`
	Data      any    `json:"data"`
}

func (e *APIException) ToJson() string {
	dj, _ := json.Marshal(e)
	return string(dj)
}

func (e *APIException) Error() string {
	return e.Message
}

// Code exception's code, 如果code不存在返回-1
func (e *APIException) ErrorCode() int {
	return int(e.ErrCode)
}

func (e *APIException) WithHttpCode(httpCode int) {
	e.HttpCode = httpCode
}

// Code exception's code, 如果code不存在返回-1
func (e *APIException) GetHttpCode() int {
	return int(e.HttpCode)
}

// WithMeta 携带一些额外信息
func (e *APIException) WithMeta(m interface{}) *APIException {
	e.Meta = m
	return e
}

func (e *APIException) GetMeta() interface{} {
	return e.Meta
}

func (e *APIException) WithData(d interface{}) *APIException {
	e.Data = d
	return e
}

func (e *APIException) GetData() interface{} {
	return e.Data
}

func (e *APIException) Is(t error) bool {
	if v, ok := t.(*APIException); ok {
		return e.ErrorCode() == v.ErrorCode()
	}

	return e.Message == t.Error()
}

func (e *APIException) GetNamespace() string {
	return e.Namespace
}

func (e *APIException) GetReason() string {
	return e.Reason
}

func (e *APIException) WithNamespace(ns string) {
	e.Namespace = ns
}
