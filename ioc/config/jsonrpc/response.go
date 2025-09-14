package jsonrpc

import "github.com/infraboard/mcube/v2/exception"

func NewResponse[T any]() *Response[T] {
	return &Response[T]{
		JSONRPC: "2.0",
		Result:  new(T),
	}
}

// JSON-RPC 2.0 响应结构
type Response[T any] struct {
	JSONRPC string                  `json:"jsonrpc"`
	Result  *T                      `json:"result,omitempty"`
	Error   *exception.ApiException `json:"error,omitempty"`
	ID      any                     `json:"id"`
}

func (r *Response[T]) SetID(id any) *Response[T] {
	r.ID = id
	return r
}
