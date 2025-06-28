package main

import (
	"github.com/infraboard/mcube/v2/examples/jsonrpc/service"
	"github.com/infraboard/mcube/v2/ioc/config/jsonrpc"
	"github.com/infraboard/mcube/v2/ioc/server/cmd"
)

var _ service.HelloService = (*HelloServiceServer)(nil)

// 实现业务功能
// req := &HelloRequest{}
// resp := &HelloResponse{}
// err := &HelloServiceServer{}.Hello(req, resp)
// net/rpc
// 1. 写好的对象， 注册给RPC Server
// 2. 再把RPC Server 启动起来
type HelloServiceServer struct {
}

// HTTP Handler
func (h *HelloServiceServer) Hello(request *service.HelloRequest, response *service.HelloResponse) error {
	response.Message = "hello:" + request.MyName
	return nil
}

func main() {
	// 注册服务
	err := jsonrpc.GetService().Registry(service.APP_NAME, &HelloServiceServer{})
	if err != nil {
		panic(err)
	}

	cmd.Start()
}
