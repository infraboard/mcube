package main

import (
	"fmt"

	"github.com/infraboard/mcube/v2/examples/jsonrpc/service"
)

func main() {
	// 1. 通过网络调用 服务端的函数(RPC)
	// 建立网络连接
	client, err := service.NewClient("http://127.0.0.1:9090/jsonrpc/mcube_app")
	if err != nil {
		panic(err)
	}
	// 方法调用
	// serviceMethod string, args any, reply any
	req := &service.HelloRequest{
		MyName: "bob",
	}
	resp := &service.HelloResponse{}
	if err := client.Hello(req, resp); err != nil {
		panic(err)
	}

	fmt.Println("xxx", resp.Message)
}
