package jsonrpc

import (
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

func Get() *JsonRpc {
	return ioc.Api().Get(APP_NAME).(*JsonRpc)
}
