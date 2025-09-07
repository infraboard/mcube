package jsonrpc

import "github.com/infraboard/mcube/v2/exception"

var (
	ErrProtocalError          = exception.NewApiException(32600, "Invalid Request").WithHttpCode(400)
	ErrMethodNotFound         = exception.NewApiException(32601, "Method not found").WithHttpCode(404)
	ErrInvalidParams          = exception.NewApiException(32602, "Invalida Params").WithHttpCode(400)
	ErrInvalidMethodSignature = exception.NewApiException(32603, "invalid method signature").WithHttpCode(400)
)
