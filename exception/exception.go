package exception

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// NewAPIException 创建一个API异常
// 用于其他模块自定义异常
func NewAPIException(namespace string, code int, reason, format string, a ...interface{}) APIException {
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

	return &exception{
		Namespace: namespace,
		Code:      code,
		Reason:    reason,
		HttpCode:  httpCode,
		Message:   fmt.Sprintf(format, a...),
	}
}

func NewAPIExceptionFromString(msg string) APIException {
	e := &exception{}
	if !strings.HasPrefix(msg, "{") {
		e.Message = msg
		e.Code = InternalServerError
		e.HttpCode = InternalServerError
		return e
	}

	err := json.Unmarshal([]byte(msg), e)
	if err != nil {
		e.Message = msg
		e.Code = InternalServerError
		e.HttpCode = InternalServerError
	}
	return e
}

// {"namespace":"","http_code":404,"error_code":404,"reason":"资源未找到","message":"test","meta":null,"data":null}
func NewAPIExceptionFromError(err error) APIException {
	return NewAPIExceptionFromString(err.Error())
}

// NewUnauthorized 未认证
func NewUnauthorized(format string, a ...interface{}) APIException {
	return NewAPIException("", Unauthorized, codeReason(Unauthorized), format, a...)
}

// NewPermissionDeny 没有权限访问
func NewPermissionDeny(format string, a ...interface{}) APIException {
	return NewAPIException("", Forbidden, codeReason(Forbidden), format, a...)
}

// NewAccessTokenIllegal 访问token过期
func NewAccessTokenIllegal(format string, a ...interface{}) APIException {
	return NewAPIException("", AccessTokenIllegal, codeReason(AccessTokenIllegal), format, a...)
}

// NewRefreshTokenIllegal 访问token过期
func NewRefreshTokenIllegal(format string, a ...interface{}) APIException {
	return NewAPIException("", RefreshTokenIllegal, codeReason(RefreshTokenIllegal), format, a...)
}

// NewOtherClientsLoggedIn 其他端登录
func NewOtherClientsLoggedIn(format string, a ...interface{}) APIException {
	return NewAPIException("", OtherClientsLoggedIn, codeReason(OtherClientsLoggedIn), format, a...)
}

// NewOtherPlaceLoggedIn 异地登录
func NewOtherPlaceLoggedIn(format string, a ...interface{}) APIException {
	return NewAPIException("", OtherPlaceLoggedIn, codeReason(OtherPlaceLoggedIn), format, a...)
}

// NewOtherIPLoggedIn 异常IP登录
func NewOtherIPLoggedIn(format string, a ...interface{}) APIException {
	return NewAPIException("", OtherIPLoggedIn, codeReason(OtherIPLoggedIn), format, a...)
}

// NewSessionTerminated 会话结束
func NewSessionTerminated(format string, a ...interface{}) APIException {
	return NewAPIException("", SessionTerminated, codeReason(SessionTerminated), format, a...)
}

// NewAccessTokenExpired 访问token过期
func NewAccessTokenExpired(format string, a ...interface{}) APIException {
	return NewAPIException("", AccessTokenExpired, codeReason(SessionTerminated), format, a...)
}

// NewRefreshTokenExpired 刷新token过期
func NewRefreshTokenExpired(format string, a ...interface{}) APIException {
	return NewAPIException("", RefreshTokenExpired, codeReason(RefreshTokenExpired), format, a...)
}

// NewBadRequest todo
func NewBadRequest(format string, a ...interface{}) APIException {
	return NewAPIException("", BadRequest, codeReason(BadRequest), format, a...)
}

// NewNotFound todo
func NewNotFound(format string, a ...interface{}) APIException {
	return NewAPIException("", NotFound, codeReason(NotFound), format, a...)
}

// NewConflict todo
func NewConflict(format string, a ...interface{}) APIException {
	return NewAPIException("", Conflict, codeReason(Conflict), format, a...)
}

// NewInternalServerError 500
func NewInternalServerError(format string, a ...interface{}) APIException {
	return NewAPIException("", InternalServerError, codeReason(InternalServerError), format, a...)
}

// NewVerifyCodeRequiredError 50018
func NewVerifyCodeRequiredError(format string, a ...interface{}) APIException {
	return NewAPIException("", VerifyCodeRequired, codeReason(VerifyCodeRequired), format, a...)
}

// NewPasswordExired 50019
func NewPasswordExired(format string, a ...interface{}) APIException {
	return NewAPIException("", PasswordExired, codeReason(PasswordExired), format, a...)
}

// IsNotFoundError 判断是否是NotFoundError
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	e, ok := err.(APIException)
	if !ok {
		return false
	}

	return e.ErrorCode() == NotFound && e.GetNamespace() == ""
}

// IsConflictError 判断是否是Conflict
func IsConflictError(err error) bool {
	if err == nil {
		return false
	}

	e, ok := err.(APIException)
	if !ok {
		return false
	}

	return e.ErrorCode() == Conflict && e.GetNamespace() == ""
}
