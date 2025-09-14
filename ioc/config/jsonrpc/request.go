package jsonrpc

// JSON-RPC 2.0 请求结构
type Request[T any] struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  T      `json:"params"`
	ID      any    `json:"id"`
}

func (r *Request[T]) SetID(id any) *Request[T] {
	r.ID = id
	return r
}

func NewRequest[T any](method string, params T) *Request[T] {
	return &Request[T]{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}
}
