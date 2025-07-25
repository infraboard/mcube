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

// NewApiException 创建一个API异常
// 用于其他模块自定义异常
func NewApiException(code int, reason string) *ApiException {
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

	return &ApiException{
		Code:     code,
		Reason:   reason,
		HttpCode: httpCode,
		Meta:     map[string]any{},
	}
}

func NewApiExceptionFromString(msg string) *ApiException {
	e := &ApiException{}
	if !strings.HasPrefix(msg, "{") {
		e.Message = msg
		e.Code = CODE_INTERNAL_SERVER_ERROR
		e.HttpCode = CODE_INTERNAL_SERVER_ERROR
		return e
	}

	err := json.Unmarshal([]byte(msg), e)
	if err != nil {
		e.Message = msg
		e.Code = CODE_INTERNAL_SERVER_ERROR
		e.HttpCode = CODE_INTERNAL_SERVER_ERROR
	}
	return e
}

func IsApiException(err error, code int) bool {
	var apiErr *ApiException
	if errors.As(err, &apiErr) {
		return apiErr.Code == code
	}
	return false
}

// ApiException API异常
type ApiException struct {
	Service  string         `json:"service"`
	HttpCode int            `json:"http_code,omitempty"`
	Code     int            `json:"code"`
	Reason   string         `json:"reason"`
	Message  string         `json:"message"`
	Meta     map[string]any `json:"meta"`
	Data     any            `json:"data"`
}

func (e *ApiException) ToJson() string {
	dj, _ := json.Marshal(e)
	return string(dj)
}

func (e *ApiException) Error() string {
	msg := e.Reason
	if e.Message != "" {
		msg += ": " + e.Message
	}
	return msg
}

// Code exception's code, 如果code不存在返回-1
func (e *ApiException) ErrorCode() int {
	return int(e.Code)
}

func (e *ApiException) WithHttpCode(httpCode int) *ApiException {
	e.HttpCode = httpCode
	return e
}

// Code exception's code, 如果code不存在返回-1
func (e *ApiException) GetHttpCode() int {
	if e.HttpCode/100 >= 1 && e.HttpCode/100 <= 5 {
		return e.HttpCode
	}
	return http.StatusInternalServerError
}

// WithMeta 携带一些额外信息
func (e *ApiException) WithMeta(key string, value any) *ApiException {
	e.Meta[key] = value
	return e
}

func (e *ApiException) GetMeta(key string) any {
	return e.Meta[key]
}

func (e *ApiException) WithData(d any) *ApiException {
	e.Data = d
	return e
}

func (e *ApiException) WithMessage(m string) *ApiException {
	e.Message = m
	return e
}

func (e *ApiException) WithMessagef(format string, a ...any) *ApiException {
	e.Message = fmt.Sprintf(format, a...)
	return e
}

func (e *ApiException) GetData() interface{} {
	return e.Data
}

func (e *ApiException) Is(t error) bool {
	if v, ok := t.(*ApiException); ok {
		return e.ErrorCode() == v.ErrorCode()
	}

	return e.Message == t.Error()
}

func (e *ApiException) GetNamespace() string {
	return e.Service
}

func (e *ApiException) GetReason() string {
	return e.Reason
}

func (e *ApiException) WithNamespace(ns string) *ApiException {
	e.Service = ns
	return e
}
