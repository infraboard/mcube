package jsonrpc

import "github.com/infraboard/mcube/v2/exception"

// JSON-RPC 2.0 响应结构
type Response struct {
	JSONRPC string                  `json:"jsonrpc"`
	Result  any                     `json:"result,omitempty"`
	Error   *exception.ApiException `json:"error,omitempty"`
	ID      any                     `json:"id"`
}
