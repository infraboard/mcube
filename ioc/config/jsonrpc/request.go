package jsonrpc

import "encoding/json"

// JSON-RPC 2.0 请求结构
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      any             `json:"id"`
}
