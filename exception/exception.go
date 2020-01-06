package exception

// NewAPIException 创建一个API异常
// 用于其他模块自定义异常
func NewAPIException(namespace string, code int, reason, format string, a ...interface{}) WithAPIException {
	// 0表示正常状态, 但是要排除变量的零值
	if code == 0 {
		code = -1
	}

	return newException(Namespace(namespace), code, reason, format, a...)
}

// NewUnauthorized 未认证
func NewUnauthorized(format string, a ...interface{}) APIException {
	return newException(usedNamespace, Unauthorized, codeReason(Unauthorized), format, a...)
}

// NewPermissionDeny 没有权限访问
func NewPermissionDeny(format string, a ...interface{}) APIException {
	return newException(usedNamespace, Forbidden, codeReason(Forbidden), format, a...)
}

// NewAccessTokenExpired 访问token过期
func NewAccessTokenExpired(format string, a ...interface{}) APIException {
	return newException(usedNamespace, AccessTokenExpired, codeReason(AccessTokenExpired), format, a...)
}

// NewRefreshTokenExpired 刷新token过期
func NewRefreshTokenExpired(format string, a ...interface{}) APIException {
	return newException(usedNamespace, RefreshTokenExpired, codeReason(RefreshTokenExpired), format, a...)
}

// NewBadRequest todo
func NewBadRequest(format string, a ...interface{}) APIException {
	return newException(usedNamespace, BadRequest, codeReason(BadRequest), format, a...)
}

// NewNotFound todo
func NewNotFound(format string, a ...interface{}) APIException {
	return newException(usedNamespace, NotFound, codeReason(BadRequest), format, a...)
}

// NewInternalServerError 500
func NewInternalServerError(format string, a ...interface{}) APIException {
	return newException(usedNamespace, InternalServerError, codeReason(InternalServerError), format, a...)
}

// IsNotFoundError 判断是否是NotFoundError
func IsNotFoundError(err error) bool {
	e, ok := err.(APIException)
	if !ok {
		return false
	}

	return e.ErrorCode() == NotFound && e.Namespace() == GlobalNamespace.String()
}
