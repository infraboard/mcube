package router_test

import (
	"net/http"
	"testing"

	"github.com/infraboard/mcube/http/router"
	"github.com/stretchr/testify/require"
)

type handler struct{}

func (h *handler) FuncWithStruct(w http.ResponseWriter, r *http.Request) {}
func FuncNoStruct(w http.ResponseWriter, r *http.Request)                {}

func TestGetHandlerFuncNameWithStruct(t *testing.T) {
	should := require.New(t)

	h := new(handler)
	fn := router.GetHandlerFuncName(h.FuncWithStruct)
	should.Equal("FuncWithStruct", fn)
}

func TestGetHandlerFuncNameWithNoStruct(t *testing.T) {
	should := require.New(t)

	fn := router.GetHandlerFuncName(FuncNoStruct)
	should.Equal("FuncNoStruct", fn)
}
