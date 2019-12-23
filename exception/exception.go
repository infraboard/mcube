package exception

// NewUnauthorized 未认证
func NewUnauthorized(format string, a ...interface{}) APIException {
	return newException(usedNamespace, Unauthorized, codeReason(Unauthorized), format, a...)
}

// NewPermissionDeny 没有权限访问
func NewPermissionDeny(format string, a ...interface{}) APIException {
	return newException(usedNamespace, Forbidden, codeReason(Forbidden), format, a...)
}

// NewTokenExpired token过期
func NewTokenExpired(format string, a ...interface{}) APIException {
	return newException(usedNamespace, TokenExpired, codeReason(TokenExpired), format, a...)
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

	return e.ErrorCode() == NotFound
}
