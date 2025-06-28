package jsonrpc

type Response[T any] struct {
	Id     uint64  `json:"id"`
	Result T       `json:"result"`
	Error  *string `json:"error"`
}

// NewResponse 创建一个新的 Response 实例
func NewResponse[T any](id uint64, result T) *Response[T] {
	return &Response[T]{
		Id:     id,
		Result: result,
	}
}
