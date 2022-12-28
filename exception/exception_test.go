package exception_test

import (
	"errors"
	"testing"

	"github.com/infraboard/mcube/exception"
)

func TestNewNotFound(t *testing.T) {
	e := exception.NewNotFound("test")
	t.Log(e.ToJson())
}

func TestNewAPIExceptionFromError(t *testing.T) {
	err := errors.New("test")

	e := exception.NewAPIExceptionFromError(err)
	t.Log(e.ToJson())
}

func TestErrLoad(t *testing.T) {
	e := exception.NewNotFound("test")

	e = exception.NewAPIExceptionFromString(e.ToJson())
	t.Log(e.ToJson())
}
