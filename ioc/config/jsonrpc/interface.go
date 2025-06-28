package jsonrpc

import (
	"io"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/ioc"
)

const (
	APP_NAME = "jsonrpc"
)

func RootRouter() *restful.Container {
	return ioc.Api().Get(APP_NAME).(*JsonRpc).Container
}

func Priority() int {
	return ioc.Config().Get(APP_NAME).Priority()
}

func GetService() Service {
	return ioc.Api().Get(APP_NAME).(Service)
}

type Service interface {
	Registry(name string, svc any) error
}

func NewRPCReadWriteCloserFromHTTP(w http.ResponseWriter, r *http.Request) *RPCReadWriteCloser {
	// 强制设置 Content-Type
	w.Header().Set("Content-Type", "application/json")
	return &RPCReadWriteCloser{w, r.Body}
}

type RPCReadWriteCloser struct {
	io.Writer
	io.ReadCloser
}

type service struct {
	name   string
	fnName string
}
