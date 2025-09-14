package main

import (
	"context"
	"fmt"

	"github.com/infraboard/mcube/v2/examples/jsonrpc/service"
	"github.com/infraboard/mcube/v2/ioc/config/jsonrpc"
	"github.com/infraboard/mcube/v2/ioc/server/cmd"
)

// 用户服务
type UserService struct{}

// RPC 方法：自动注册为 getUser
func (s *UserService) RPCGetUser(ctx context.Context, req *service.HelloRequest) (*service.HelloResponse, error) {
	// 直接使用解析好的请求对象
	return &service.HelloResponse{
		Message: fmt.Sprintf("Hello, %s", req.MyName),
	}, nil
}

func main() {
	jsonrpc.RegisterService(&UserService{})
	cmd.Start()
}
