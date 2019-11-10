package exception

import "net/http"

const (
	TokenExpired = 1000

	// 1xx - 5xx copy from http status code
	Unauthorized        = http.StatusUnauthorized
	BadRequest          = http.StatusBadRequest
	InternalServerError = http.StatusInternalServerError
	Forbidden           = http.StatusForbidden
	NotFound            = http.StatusNotFound

	UnKnownException = 9999
)
