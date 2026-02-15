package exception

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// Unwrap 返回包装的原始错误，支持 Go 1.13+ 错误链
func (e *ApiException) Unwrap() error {
	return e.cause
}

// Cause 返回根本原因错误
func (e *ApiException) Cause() error {
	return e.cause
}

// WithCause 设置原始错误
func (e *ApiException) WithCause(cause error) *ApiException {
	e.cause = cause
	return e
}

// Wrap 包装一个错误，支持错误链追踪
// 这是推荐的创建异常的方式，可以保留原始错误信息
func Wrap(err error, code int, reason string) *ApiException {
	if err == nil {
		return nil
	}

	// 如果已经是ApiException，直接返回
	var apiErr *ApiException
	if errors.As(err, &apiErr) {
		return apiErr
	}

	e := NewApiException(code, reason)
	e.cause = err
	e.Message = err.Error()
	return e
}

// Wrapf 包装错误并格式化消息
func Wrapf(err error, code int, reason string, format string, args ...any) *ApiException {
	if err == nil {
		return nil
	}

	e := Wrap(err, code, reason)
	e.Message = fmt.Sprintf(format, args...)
	return e
}

// ErrorChain 返回完整的错误链信息
func (e *ApiException) ErrorChain() []string {
	var chain []string
	chain = append(chain, e.Error())

	err := e.cause
	for err != nil {
		chain = append(chain, err.Error())
		if unwrapper, ok := err.(interface{ Unwrap() error }); ok {
			err = unwrapper.Unwrap()
		} else {
			break
		}
	}
	return chain
}

// ErrorChainString 返回格式化的错误链字符串
func (e *ApiException) ErrorChainString() string {
	chain := e.ErrorChain()
	if len(chain) <= 1 {
		return e.Error()
	}
	return strings.Join(chain, " → ")
}

// WithStack 添加堆栈追踪信息（默认3层）
func (e *ApiException) WithStack() *ApiException {
	return e.WithStackDepth(3)
}

// WithStackDepth 添加指定层数的堆栈追踪
func (e *ApiException) WithStackDepth(depth int) *ApiException {
	if depth <= 0 {
		depth = 3
	}
	e.stack = captureStack(depth, 10) // 最多10层
	return e
}

// GetStack 获取堆栈信息
func (e *ApiException) GetStack() string {
	return e.stack
}

// captureStack 捕获堆栈信息
func captureStack(skip, maxDepth int) string {
	var buf strings.Builder
	for i := skip; i < skip+maxDepth; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			break
		}
		// 简化文件路径，只显示相对路径
		shortFile := file
		if idx := strings.LastIndex(file, "/"); idx >= 0 {
			shortFile = file[idx+1:]
		}
		buf.WriteString(fmt.Sprintf("\n  at %s (%s:%d)", fn.Name(), shortFile, line))
	}
	return buf.String()
}

// WithTraceID 设置追踪ID
func (e *ApiException) WithTraceID(traceID string) *ApiException {
	e.TraceID = traceID
	e.safeWithMeta("trace_id", traceID)
	return e
}

// GetTraceID 获取追踪ID
func (e *ApiException) GetTraceID() string {
	return e.TraceID
}

// WithRequestID 设置请求ID
func (e *ApiException) WithRequestID(requestID string) *ApiException {
	e.RequestID = requestID
	return e
}

// GetRequestID 获取请求ID
func (e *ApiException) GetRequestID() string {
	return e.RequestID
}

// 线程安全的Meta操作
func (e *ApiException) safeWithMeta(key string, value any) *ApiException {
	e.metaMu.Lock()
	defer e.metaMu.Unlock()
	if e.Meta == nil {
		e.Meta = make(map[string]any)
	}
	e.Meta[key] = value
	return e
}

func (e *ApiException) safeGetMeta(key string) any {
	e.metaMu.RLock()
	defer e.metaMu.RUnlock()
	return e.Meta[key]
}

// ErrorType 错误类型分类
type ErrorType int

const (
	ErrorTypeUnknown    ErrorType = iota
	ErrorTypeClient               // 客户端错误（4xx）
	ErrorTypeServer               // 服务端错误（5xx）
	ErrorTypeAuth                 // 认证授权错误
	ErrorTypeValidation           // 验证错误
	ErrorTypeNotFound             // 资源不存在
	ErrorTypeConflict             // 资源冲突
	ErrorTypeTimeout              // 超时错误
	ErrorTypeRateLimit            // 限流错误
)

// Type 返回错误类型
func (e *ApiException) Type() ErrorType {
	switch e.Code {
	case CODE_UNAUTHORIZED, CODE_ACCESS_TOKEN_ILLEGAL, CODE_REFRESH_TOKEN_ILLEGAL,
		CODE_ACESS_TOKEN_EXPIRED, CODE_REFRESH_TOKEN_EXPIRED:
		return ErrorTypeAuth
	case CODE_BAD_REQUEST, CODE_VERIFY_CODE_REQUIRED:
		return ErrorTypeValidation
	case CODE_NOT_FOUND:
		return ErrorTypeNotFound
	case CODE_CONFLICT:
		return ErrorTypeConflict
	case CODE_FORBIDDEN:
		return ErrorTypeAuth
	case CODE_INTERNAL_SERVER_ERROR:
		return ErrorTypeServer
	default:
		// 根据HTTP状态码判断
		if e.Code >= 400 && e.Code < 500 {
			return ErrorTypeClient
		} else if e.Code >= 500 {
			return ErrorTypeServer
		}
		return ErrorTypeUnknown
	}
}

// IsRetryable 判断错误是否可重试
func (e *ApiException) IsRetryable() bool {
	// 服务端错误（5xx）通常可以重试
	if e.Type() == ErrorTypeServer {
		return true
	}
	// 超时和限流错误也可以重试
	if e.Type() == ErrorTypeTimeout || e.Type() == ErrorTypeRateLimit {
		return true
	}
	// 某些特定的错误也可以重试
	switch e.Code {
	case CODE_SESSION_TERMINATED, CODE_ACESS_TOKEN_EXPIRED:
		return true
	}
	return false
}

// IsTemporary 判断是否是临时错误
func (e *ApiException) IsTemporary() bool {
	return e.IsRetryable()
}

// IsClientError 判断是否是客户端错误（4xx）
func (e *ApiException) IsClientError() bool {
	return e.Type() == ErrorTypeClient || e.Type() == ErrorTypeValidation ||
		e.Type() == ErrorTypeAuth || e.Type() == ErrorTypeNotFound ||
		e.Type() == ErrorTypeConflict
}

// IsServerError 判断是否是服务端错误（5xx）
func (e *ApiException) IsServerError() bool {
	return e.Type() == ErrorTypeServer
}

// IsAuthError 判断是否是认证/授权错误
func (e *ApiException) IsAuthError() bool {
	return e.Type() == ErrorTypeAuth
}
