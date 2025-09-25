package jsonrpc

import (
	"context"
	"net/http"
)

type Auther interface {
	// RPC认证接口, 返回认证信息, 携带在RpcContext中, 可以通过GetRpcContext获取
	Auth(context.Context, *AuthRequest) (authInfo any, err error)
}

type AuthRequest struct {
	Header *http.Header `json:"header"`
	Method string       `json:"method"`
}

func SetAuther(a Auther) {
	Get().auther = a
}
