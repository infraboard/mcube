package exception_test

import (
	"errors"
	"testing"

	"github.com/infraboard/mcube/exception"
)

func TestNewNotFound(t *testing.T) {
	e := exception.NewNotFound("test")
	t.Log(e)

	e = exception.NewAPIExceptionFromError(e)
	t.Log(e)
}

func TestNewAPIExceptionFromError(t *testing.T) {
	e := errors.New("test")

	e = exception.NewAPIExceptionFromError(e)
	t.Log(e)
}
