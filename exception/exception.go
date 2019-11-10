package exception

import "fmt"

// NewUnauthorized 未认证
func NewUnauthorized(format string, a ...interface{}) APIException {
	return &exception{
		msg:  fmt.Sprintf(format, a...),
		code: Unauthorized,
	}
}

// NewPermissionDeny 没有权限访问
func NewPermissionDeny(format string, a ...interface{}) APIException {
	return &exception{
		msg:  fmt.Sprintf(format, a...),
		code: Forbidden,
	}
}

// NewBadRequest todo
func NewBadRequest(format string, a ...interface{}) APIException {
	return &exception{
		msg:  fmt.Sprintf(format, a...),
		code: BadRequest,
	}
}

// NewNotFound todo
func NewNotFound(format string, a ...interface{}) APIException {
	return &exception{
		msg:  fmt.Sprintf(format, a...),
		code: NotFound,
	}
}

// NewInternalServerError 500
func NewInternalServerError(format string, a ...interface{}) APIException {
	return &exception{
		msg:  fmt.Sprintf(format, a...),
		code: InternalServerError,
	}
}

// IsNotFoundError 判断是否是NotFoundError
func IsNotFoundError(err error) bool {
	e, ok := err.(APIException)
	if !ok {
		return false
	}

	return e.ErrorCode() == NotFound
}
