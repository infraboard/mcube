package jsonrpc

func NewRequest[T any](method string, Param T) *Request[T] {
	return &Request[T]{
		Method: method,
		Params: [1]T{Param},
	}
}

type Request[T any] struct {
	Method string `json:"method"`
	Params [1]T   `json:"params"`
	Id     uint64 `json:"id"`
}
