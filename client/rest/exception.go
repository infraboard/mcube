package rest

import (
	"fmt"

	"github.com/infraboard/mcube/v2/client/negotiator"
)

type ExceptionHandleFunc func(*Exception) error

func NewException(code int, body []byte) *Exception {
	return &Exception{
		Code: code,
		Body: body,
	}
}

type Exception struct {
	Code    int
	Body    []byte
	decoder negotiator.Decoder
}

func (e *Exception) WithDecoder(decoder negotiator.Decoder) *Exception {
	e.decoder = decoder
	return e
}

func (e *Exception) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", e.Code, string(e.Body))
}

func (e *Exception) Decode(v any) error {
	if e.decoder == nil {
		return fmt.Errorf("with decoder required")
	}

	return e.decoder.Decode(e.Body, v)
}

func (e *Exception) MustDecode(v any) {
	if err := e.Decode(v); err != nil {
		panic(err)
	}
}
