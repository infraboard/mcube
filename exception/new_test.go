package exception_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/infraboard/mcube/v2/exception"
)

// 测试错误链支持
func TestWrap(t *testing.T) {
	originalErr := errors.New("database connection failed")

	wrapped := exception.Wrap(originalErr, exception.CODE_INTERNAL_SERVER_ERROR, "数据库错误")

	if wrapped == nil {
		t.Fatal("wrapped error should not be nil")
	}

	// 测试Unwrap
	if errors.Unwrap(wrapped) != originalErr {
		t.Error("Unwrap should return original error")
	}

	// 测试Is
	if !errors.Is(wrapped, originalErr) {
		t.Error("errors.Is should work with wrapped errors")
	}

	// 测试错误链
	chain := wrapped.ErrorChain()
	if len(chain) != 2 {
		t.Errorf("expected 2 errors in chain, got %d", len(chain))
	}

	t.Logf("Error chain: %s", wrapped.ErrorChainString())
}

func TestWithStack(t *testing.T) {
	err := exception.NewNotFound("resource not found").WithStack()

	stack := err.GetStack()
	if stack == "" {
		t.Error("stack should not be empty")
	}

	t.Logf("Stack trace: %s", stack)
}

func TestWithTraceID(t *testing.T) {
	err := exception.NewBadRequest("invalid input").
		WithTraceID("trace-123").
		WithRequestID("req-456")

	if err.GetTraceID() != "trace-123" {
		t.Error("trace ID mismatch")
	}

	if err.GetRequestID() != "req-456" {
		t.Error("request ID mismatch")
	}

	// 验证 RequestID 是直接字段
	if err.RequestID != "req-456" {
		t.Error("RequestID field should be directly accessible")
	}

	// 验证 JSON 序列化包含 RequestID
	jsonStr := err.ToJson()
	if !strings.Contains(jsonStr, "\"request_id\":\"req-456\"") {
		t.Errorf("JSON should contain request_id field: %s", jsonStr)
	}

	t.Logf("Error with IDs: %s", err.ToJson())
}

func TestErrorType(t *testing.T) {
	tests := []struct {
		name     string
		err      *exception.ApiException
		wantType exception.ErrorType
	}{
		{
			name:     "认证错误",
			err:      exception.NewUnauthorized("未登录"),
			wantType: exception.ErrorTypeAuth,
		},
		{
			name:     "验证错误",
			err:      exception.NewBadRequest("参数错误"),
			wantType: exception.ErrorTypeValidation,
		},
		{
			name:     "未找到",
			err:      exception.NewNotFound("资源不存在"),
			wantType: exception.ErrorTypeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Type(); got != tt.wantType {
				t.Errorf("Type() = %v, want %v", got, tt.wantType)
			}
		})
	}
}

func TestIsRetryable(t *testing.T) {
	serverErr := exception.NewInternalServerError("系统错误")
	if !serverErr.IsRetryable() {
		t.Error("Server error should be retryable")
	}

	clientErr := exception.NewBadRequest("参数错误")
	if clientErr.IsRetryable() {
		t.Error("Client error should not be retryable")
	}
}

func TestBuilder(t *testing.T) {
	err := exception.NewBuilder(exception.CODE_BAD_REQUEST).
		WithReason("参数验证失败").
		WithMessage("邮箱格式不正确").
		WithMeta("field", "email").
		WithTraceID("trace-789").
		WithRequestID("req-abc123").
		WithStack().
		Build()

	if err.ErrorCode() != exception.CODE_BAD_REQUEST {
		t.Error("code mismatch")
	}

	if err.GetReason() != "参数验证失败" {
		t.Error("reason mismatch")
	}

	if err.GetMeta("field") != "email" {
		t.Error("meta mismatch")
	}

	if err.GetTraceID() != "trace-789" {
		t.Error("trace ID mismatch")
	}

	if err.GetRequestID() != "req-abc123" {
		t.Error("request ID mismatch")
	}

	// 验证 RequestID 字段可直接访问
	if err.RequestID != "req-abc123" {
		t.Error("RequestID field should be directly accessible")
	}

	t.Logf("Built error: %s", err.ToJson())
}

func TestBackwardCompatibility(t *testing.T) {
	err := exception.NewNotFound("user %s not found", "alice")
	err.WithMeta("user_id", "123")

	if err.ErrorCode() != exception.CODE_NOT_FOUND {
		t.Error("code mismatch")
	}

	if !exception.IsNotFoundError(err) {
		t.Error("IsNotFoundError should work")
	}

	t.Logf("Backward compatibility test passed: %s", err.ToJson())
}
