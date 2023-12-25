package main

import (
	"context"
	"fmt"

	"github.com/infraboard/mcube/v2/examples/rpc_grpc/pb"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/server"
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

func (h *HelloGrpc) Greet(ctx context.Context, in *pb.GreetRequest) (*pb.GreetResponse, error) {
	return &pb.GreetResponse{
		Msg: fmt.Sprintf("hello, %s", in.Name),
	}, nil
}

func main() {
	// 注册HTTP接口类
	ioc.Controller().Registry(&HelloGrpc{})

	// 启动应用
	err := server.Run(context.Background())
	if err != nil {
		panic(err)
	}
}
