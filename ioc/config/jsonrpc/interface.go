package jsonrpc

import (
	"io"
	"net/http"

	"github.com/infraboard/mcube/v2/ioc"
)

const (
	APP_NAME = "jsonrpc"
)

func GetService() Service {
	return ioc.Api().Get(APP_NAME).(Service)
}

type Service interface {
	Registry(name string, svc any) error
}

func NewRPCReadWriteCloserFromHTTP(w http.ResponseWriter, r *http.Request) *RPCReadWriteCloser {
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
