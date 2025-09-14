package jsonrpc

import (
	"context"
	"net/http"
)

type Auther interface {
	Auth(context.Context, *AuthRequest) (*RpcContext, error)
}

type AuthRequest struct {
	Header *http.Header `json:"header"`
	Method string       `json:"method"`
}
