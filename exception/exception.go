package exception

import (
	"encoding/json"
	"errors"
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
func NewAPIException(code int, reason string) *APIException {
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
		BizCode:  code,
		Reason:   reason,
		HttpCode: httpCode,
	}
}

func NewAPIExceptionFromString(msg string) *APIException {
	e := &APIException{}
	if !strings.HasPrefix(msg, "{") {
		e.Message = msg
		e.BizCode = InternalServerError
		e.HttpCode = InternalServerError
		return e
	}

	err := json.Unmarshal([]byte(msg), e)
	if err != nil {
		e.Message = msg
		e.BizCode = InternalServerError
		e.HttpCode = InternalServerError
	}
	return e
}

func IsError(err error, targetError *APIException) bool {
	var apiErr *APIException
	if errors.As(err, &apiErr) {
		return apiErr.BizCode == targetError.BizCode
	}
	return false
}

func IsAPIException(err error, code int) bool {
	var apiErr *APIException
	if errors.As(err, &apiErr) {
		return apiErr.BizCode == code
	}
	return false
}

// APIException API异常
type APIException struct {
	Namespace string `json:"namespace"`
	HttpCode  int    `json:"http_code,omitempty"`
	BizCode   int    `json:"code"`
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
	return int(e.BizCode)
}

func (e *APIException) WithHttpCode(httpCode int) *APIException {
	e.HttpCode = httpCode
	return e
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

func (e *APIException) WithMessage(m string) *APIException {
	e.Message = m
	return e
}

func (e *APIException) WithMessagef(format string, a ...any) *APIException {
	e.Message = fmt.Sprintf(format, a...)
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

func (e *APIException) WithNamespace(ns string) *APIException {
	e.Namespace = ns
	return e
}
