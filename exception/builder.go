package exception

import "fmt"

// ExceptionBuilder 异常构建器，提供更灵活的异常创建方式
type ExceptionBuilder struct {
	code      int
	reason    string
	message   string
	httpCode  int
	meta      map[string]any
	data      any
	cause     error
	service   string
	traceID   string
	requestID string
	withStack bool
}

// NewBuilder 创建异常构建器
func NewBuilder(code int) *ExceptionBuilder {
	return &ExceptionBuilder{
		code:     code,
		reason:   codeReason(code),
		httpCode: code,
		meta:     make(map[string]any),
	}
}

// WithReason 设置错误原因
func (b *ExceptionBuilder) WithReason(reason string) *ExceptionBuilder {
	b.reason = reason
	return b
}

// WithMessage 设置错误消息
func (b *ExceptionBuilder) WithMessage(message string) *ExceptionBuilder {
	b.message = message
	return b
}

// WithMessagef 设置格式化的错误消息
func (b *ExceptionBuilder) WithMessagef(format string, args ...any) *ExceptionBuilder {
	b.message = fmt.Sprintf(format, args...)
	return b
}

// WithHTTPCode 设置HTTP状态码
func (b *ExceptionBuilder) WithHTTPCode(httpCode int) *ExceptionBuilder {
	b.httpCode = httpCode
	return b
}

// WithMeta 添加元数据
func (b *ExceptionBuilder) WithMeta(key string, value any) *ExceptionBuilder {
	b.meta[key] = value
	return b
}

// WithData 设置响应数据
func (b *ExceptionBuilder) WithData(data any) *ExceptionBuilder {
	b.data = data
	return b
}

// WithCause 设置原始错误
func (b *ExceptionBuilder) WithCause(cause error) *ExceptionBuilder {
	b.cause = cause
	if b.message == "" && cause != nil {
		b.message = cause.Error()
	}
	return b
}

// WithService 设置服务名称
func (b *ExceptionBuilder) WithService(service string) *ExceptionBuilder {
	b.service = service
	return b
}

// WithTraceID 设置追踪ID
func (b *ExceptionBuilder) WithTraceID(traceID string) *ExceptionBuilder {
	b.traceID = traceID
	return b
}

// WithRequestID 设置请求ID
func (b *ExceptionBuilder) WithRequestID(requestID string) *ExceptionBuilder {
	b.requestID = requestID
	return b
}

// WithStack 启用堆栈追踪
func (b *ExceptionBuilder) WithStack() *ExceptionBuilder {
	b.withStack = true
	return b
}

// Build 构建API异常
func (b *ExceptionBuilder) Build() *ApiException {
	e := NewApiException(b.code, b.reason)

	if b.message != "" {
		e.Message = b.message
	}

	if b.httpCode > 0 {
		e.HttpCode = b.httpCode
	}

	if len(b.meta) > 0 {
		e.Meta = b.meta
	}

	if b.data != nil {
		e.Data = b.data
	}

	if b.cause != nil {
		e.cause = b.cause
	}

	if b.service != "" {
		e.Service = b.service
	}

	if b.traceID != "" {
		e.TraceID = b.traceID
	}

	if b.requestID != "" {
		e.RequestID = b.requestID
	}

	if b.withStack {
		e.WithStack()
	}

	return e
}
