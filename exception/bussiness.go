package exception

// {"namespace":"","http_code":404,"error_code":404,"reason":"资源未找到","message":"test","meta":null,"data":null}
func NewAPIExceptionFromError(err error) *ApiException {
	return NewApiExceptionFromString(err.Error())
}

// NewUnauthorized 未认证
func NewUnauthorized(format string, a ...interface{}) *ApiException {
	return NewApiException(CODE_UNAUTHORIZED, codeReason(CODE_UNAUTHORIZED)).WithMessagef(format, a...)
}

// NewPermissionDeny 没有权限访问
func NewPermissionDeny(format string, a ...interface{}) *ApiException {
	return NewApiException(CODE_FORBIDDEN, codeReason(CODE_FORBIDDEN)).WithMessagef(format, a...)
}

// NewAccessTokenIllegal 访问token过期
func NewAccessTokenIllegal(format string, a ...interface{}) *ApiException {
	return NewApiException(CODE_ACCESS_TOKEN_ILLEGAL, codeReason(CODE_ACCESS_TOKEN_ILLEGAL)).WithMessagef(format, a...)
}

// NewRefreshTokenIllegal 访问token过期
func NewRefreshTokenIllegal(format string, a ...interface{}) *ApiException {
	return NewApiException(CODE_REFRESH_TOKEN_ILLEGAL, codeReason(CODE_REFRESH_TOKEN_ILLEGAL)).WithMessagef(format, a...)
}

// NewOtherClientsLoggedIn 其他端登录
func NewOtherClientsLoggedIn(format string, a ...interface{}) *ApiException {
	return NewApiException(CODE_OTHER_CLIENT_LOGIN, codeReason(CODE_OTHER_CLIENT_LOGIN)).WithMessagef(format, a...)
}

// NewOtherPlaceLoggedIn 异地登录
func NewOtherPlaceLoggedIn(format string, a ...interface{}) *ApiException {
	return NewApiException(CODE_OTHER_PLACE_LGOIN, codeReason(CODE_OTHER_PLACE_LGOIN)).WithMessagef(format, a...)
}

// NewOtherIPLoggedIn 异常IP登录
func NewOtherIPLoggedIn(format string, a ...interface{}) *ApiException {
	return NewApiException(CODE_OTHER_IP_LOGIN, codeReason(CODE_OTHER_IP_LOGIN)).WithMessagef(format, a...)
}

// NewSessionTerminated 会话结束
func NewSessionTerminated(format string, a ...interface{}) *ApiException {
	return NewApiException(CODE_SESSION_TERMINATED, codeReason(CODE_SESSION_TERMINATED)).WithMessagef(format, a...)
}

// NewAccessTokenExpired 访问token过期
func NewAccessTokenExpired(format string, a ...interface{}) *ApiException {
	return NewApiException(CODE_ACESS_TOKEN_EXPIRED, codeReason(CODE_SESSION_TERMINATED)).WithMessagef(format, a...)
}

// NewRefreshTokenExpired 刷新token过期
func NewRefreshTokenExpired(format string, a ...interface{}) *ApiException {
	return NewApiException(CODE_REFRESH_TOKEN_EXPIRED, codeReason(CODE_REFRESH_TOKEN_EXPIRED)).WithMessagef(format, a...)
}

// NewBadRequest todo
func NewBadRequest(format string, a ...interface{}) *ApiException {
	return NewApiException(CODE_BAD_REQUEST, codeReason(CODE_BAD_REQUEST)).WithMessagef(format, a...)
}

// NewNotFound todo
func NewNotFound(format string, a ...interface{}) *ApiException {
	return NewApiException(CODE_NOT_FOUND, codeReason(CODE_NOT_FOUND)).WithMessagef(format, a...)
}

// NewConflict todo
func NewConflict(format string, a ...interface{}) *ApiException {
	return NewApiException(CODE_CONFLICT, codeReason(CODE_CONFLICT)).WithMessagef(format, a...)
}

// NewInternalServerError 500
func NewInternalServerError(format string, a ...interface{}) *ApiException {
	return NewApiException(CODE_INTERNAL_SERVER_ERROR, codeReason(CODE_INTERNAL_SERVER_ERROR)).WithMessagef(format, a...)
}

// NewVerifyCodeRequiredError 50018
func NewVerifyCodeRequiredError(format string, a ...interface{}) *ApiException {
	return NewApiException(CODE_VERIFY_CODE_REQUIRED, codeReason(CODE_VERIFY_CODE_REQUIRED)).WithMessagef(format, a...)
}

// NewPasswordExired 50019
func NewPasswordExired(format string, a ...interface{}) *ApiException {
	return NewApiException(CODE_PASSWORD_EXPIRED, codeReason(CODE_PASSWORD_EXPIRED)).WithMessagef(format, a...)
}

// IsNotFoundError 判断是否是NotFoundError
func IsNotFoundError(err error) bool {
	return IsApiException(err, CODE_NOT_FOUND)
}

// IsConflictError 判断是否是Conflict
func IsConflictError(err error) bool {
	return IsApiException(err, CODE_CONFLICT)
}
