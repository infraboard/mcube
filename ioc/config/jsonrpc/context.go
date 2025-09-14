package jsonrpc

import "context"

type RpcContextKey struct{}

type RpcContext struct {
	AuthInfo any
}

func GetRpcContext(ctx context.Context) *RpcContext {
	if ctx == nil {
		return nil
	}
	v := ctx.Value(RpcContextKey{})
	if v == nil {
		return nil
	}
	return v.(*RpcContext)
}
