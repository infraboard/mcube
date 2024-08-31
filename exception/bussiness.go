package exception

// {"namespace":"","http_code":404,"error_code":404,"reason":"资源未找到","message":"test","meta":null,"data":null}
func NewAPIExceptionFromError(err error) *APIException {
	return NewAPIExceptionFromString(err.Error())
}

// NewUnauthorized 未认证
func NewUnauthorized(format string, a ...interface{}) *APIException {
	return NewAPIException(Unauthorized, codeReason(Unauthorized)).WithMessagef(format, a...)
}

// NewPermissionDeny 没有权限访问
func NewPermissionDeny(format string, a ...interface{}) *APIException {
	return NewAPIException(Forbidden, codeReason(Forbidden)).WithMessagef(format, a...)
}

// NewAccessTokenIllegal 访问token过期
func NewAccessTokenIllegal(format string, a ...interface{}) *APIException {
	return NewAPIException(AccessTokenIllegal, codeReason(AccessTokenIllegal)).WithMessagef(format, a...)
}

// NewRefreshTokenIllegal 访问token过期
func NewRefreshTokenIllegal(format string, a ...interface{}) *APIException {
	return NewAPIException(RefreshTokenIllegal, codeReason(RefreshTokenIllegal)).WithMessagef(format, a...)
}

// NewOtherClientsLoggedIn 其他端登录
func NewOtherClientsLoggedIn(format string, a ...interface{}) *APIException {
	return NewAPIException(OtherClientsLoggedIn, codeReason(OtherClientsLoggedIn)).WithMessagef(format, a...)
}

// NewOtherPlaceLoggedIn 异地登录
func NewOtherPlaceLoggedIn(format string, a ...interface{}) *APIException {
	return NewAPIException(OtherPlaceLoggedIn, codeReason(OtherPlaceLoggedIn)).WithMessagef(format, a...)
}

// NewOtherIPLoggedIn 异常IP登录
func NewOtherIPLoggedIn(format string, a ...interface{}) *APIException {
	return NewAPIException(OtherIPLoggedIn, codeReason(OtherIPLoggedIn)).WithMessagef(format, a...)
}

// NewSessionTerminated 会话结束
func NewSessionTerminated(format string, a ...interface{}) *APIException {
	return NewAPIException(SessionTerminated, codeReason(SessionTerminated)).WithMessagef(format, a...)
}

// NewAccessTokenExpired 访问token过期
func NewAccessTokenExpired(format string, a ...interface{}) *APIException {
	return NewAPIException(AccessTokenExpired, codeReason(SessionTerminated)).WithMessagef(format, a...)
}

// NewRefreshTokenExpired 刷新token过期
func NewRefreshTokenExpired(format string, a ...interface{}) *APIException {
	return NewAPIException(RefreshTokenExpired, codeReason(RefreshTokenExpired)).WithMessagef(format, a...)
}

// NewBadRequest todo
func NewBadRequest(format string, a ...interface{}) *APIException {
	return NewAPIException(BadRequest, codeReason(BadRequest)).WithMessagef(format, a...)
}

// NewNotFound todo
func NewNotFound(format string, a ...interface{}) *APIException {
	return NewAPIException(NotFound, codeReason(NotFound)).WithMessagef(format, a...)
}

// NewConflict todo
func NewConflict(format string, a ...interface{}) *APIException {
	return NewAPIException(Conflict, codeReason(Conflict)).WithMessagef(format, a...)
}

// NewInternalServerError 500
func NewInternalServerError(format string, a ...interface{}) *APIException {
	return NewAPIException(InternalServerError, codeReason(InternalServerError)).WithMessagef(format, a...)
}

// NewVerifyCodeRequiredError 50018
func NewVerifyCodeRequiredError(format string, a ...interface{}) *APIException {
	return NewAPIException(VerifyCodeRequired, codeReason(VerifyCodeRequired)).WithMessagef(format, a...)
}

// NewPasswordExired 50019
func NewPasswordExired(format string, a ...interface{}) *APIException {
	return NewAPIException(PasswordExired, codeReason(PasswordExired)).WithMessagef(format, a...)
}

// IsNotFoundError 判断是否是NotFoundError
func IsNotFoundError(err error) bool {
	return IsAPIException(err, NotFound)
}

// IsConflictError 判断是否是Conflict
func IsConflictError(err error) bool {
	return IsAPIException(err, Conflict)
}
