package rest

import "encoding/json"

func NewResponse() *Response {
	return &Response{}
}

// Response contains the result of calling Request.Do().
type Response struct {
	body        []byte
	contentType string
	err         error
	statusCode  int
}

func (r *Response) Into(v any) error {
	return json.Unmarshal(r.body, v)
}
