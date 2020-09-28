package exception

import "fmt"

// NewAPIException 创建一个API异常
// 用于其他模块自定义异常
func NewAPIException(namespace string, code int, reason, format string, a ...interface{}) APIException {
	// 0表示正常状态, 但是要排除变量的零值
	if code == 0 {
		code = -1
	}

	return &exception{
		namespace: Namespace(namespace),
		code:      code,
		reason:    codeReason(code),
		message:   fmt.Sprintf(format, a...),
	}
}

// NewUnauthorized 未认证
func NewUnauthorized(format string, a ...interface{}) APIException {
	return newException(usedNamespace, Unauthorized, format, a...)
}

// NewPermissionDeny 没有权限访问
func NewPermissionDeny(format string, a ...interface{}) APIException {
	return newException(usedNamespace, Forbidden, format, a...)
}

// NewAccessTokenIllegal 访问token过期
func NewAccessTokenIllegal(format string, a ...interface{}) APIException {
	return newException(usedNamespace, AccessTokenIllegal, format, a...)
}

// NewRefreshTokenIllegal 访问token过期
func NewRefreshTokenIllegal(format string, a ...interface{}) APIException {
	return newException(usedNamespace, RefreshTokenIllegal, format, a...)
}

// NewOtherClientsLoggedIn 访问token过期
func NewOtherClientsLoggedIn(format string, a ...interface{}) APIException {
	return newException(usedNamespace, OtherClientsLoggedIn, format, a...)
}

// NewAccessTokenExpired 访问token过期
func NewAccessTokenExpired(format string, a ...interface{}) APIException {
	return newException(usedNamespace, AccessTokenExpired, format, a...)
}

// NewRefreshTokenExpired 刷新token过期
func NewRefreshTokenExpired(format string, a ...interface{}) APIException {
	return newException(usedNamespace, RefreshTokenExpired, format, a...)
}

// NewBadRequest todo
func NewBadRequest(format string, a ...interface{}) APIException {
	return newException(usedNamespace, BadRequest, format, a...)
}

// NewNotFound todo
func NewNotFound(format string, a ...interface{}) APIException {
	return newException(usedNamespace, NotFound, format, a...)
}

// NewConflict todo
func NewConflict(format string, a ...interface{}) APIException {
	return newException(usedNamespace, Conflict, format, a...)
}

// NewInternalServerError 500
func NewInternalServerError(format string, a ...interface{}) APIException {
	return newException(usedNamespace, InternalServerError, format, a...)
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

	return e.ErrorCode() == NotFound && e.Namespace() == GlobalNamespace.String()
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

	return e.ErrorCode() == Conflict && e.Namespace() == GlobalNamespace.String()
}
