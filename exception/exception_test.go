package exception_test

import (
	"testing"

	"github.com/infraboard/mcube/exception"
)

func TestNewNotFound(t *testing.T) {
	e := exception.NewNotFound("test")
	t.Log(e)
}
