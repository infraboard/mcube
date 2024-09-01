package exception_test

import (
	"errors"
	"testing"

	"github.com/infraboard/mcube/v2/exception"
)

func TestNewNotFound(t *testing.T) {
	e := exception.NewNotFound("user %s not found", "alice")
	// {"service":"","http_code":404,"code":404,"reason":"资源未找到","message":"user alice not found","meta":null,"data":null}
	t.Log(e.ToJson())
	t.Log(exception.IsApiException(e, exception.CODE_NOT_FOUND))
}

func TestNewAPIExceptionFromError(t *testing.T) {
	err := errors.New("test")

	e := exception.NewAPIExceptionFromError(err)
	t.Log(e.ToJson())
}

func TestErrLoad(t *testing.T) {
	e := exception.NewNotFound("test")

	e = exception.NewApiExceptionFromString(e.ToJson())
	t.Log(e.ToJson())
}
