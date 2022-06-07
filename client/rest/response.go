package rest

import (
	"io"
	"io/ioutil"
	"net/http"
)

func NewResponse() *Response {
	return &Response{}
}

// Response contains the result of calling Request.Do().
type Response struct {
	body       io.ReadCloser
	headers    http.Header
	statusCode int
	err        error
	bf         []byte
	isRead     bool
	decoder    Decoder
}

func (r *Response) read() {
	if r.body == nil || r.isRead {
		return
	}

	r.isRead = true
	body, err := ioutil.ReadAll(r.body)
	if err != nil {
		r.err = err
		return
	}
	defer r.body.Close()

	r.bf = body
}

func (r *Response) Into(v any) error {
	if r.err != nil {
		return r.err
	}

	r.read()
	return r.decoder.Decode(r.bf, v)
}
