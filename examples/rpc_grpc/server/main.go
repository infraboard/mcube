package main

import (
	"context"

	"github.com/infraboard/mcube/v2/examples/rpc_grpc/pb"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	"google.golang.org/grpc"
)

type HelloGrpc struct {
	// 继承自Ioc对象
	ioc.ObjectImpl
	// 集成Grpc Server对象
	pb.UnimplementedHelloServer
}

func (h *HelloGrpc) Name() string {
	return "hello_module"
}

func (h *HelloGrpc) Registry(server *grpc.Server) {
	pb.RegisterHelloServer(server, h)
}

func main() {
	// 注册HTTP接口类
	ioc.Controller().Registry(&HelloGrpc{})

	// Ioc配置
	req := ioc.NewLoadConfigRequest()
	err := ioc.ConfigIocObject(req)
	if err != nil {
		panic(err)
	}

	// 启动应用
	err = application.App().Start(context.Background())
	if err != nil {
		panic(err)
	}
}
