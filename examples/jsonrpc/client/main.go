package main

import (
	"context"
	"fmt"

	"github.com/infraboard/mcube/v2/examples/jsonrpc/service"
)

func main() {
	// 1. 通过网络调用 服务端的函数(RPC)
	client, err := service.NewClient("http://127.0.0.1:9090/jsonrpc/mcube_app/v1")
	if err != nil {
		panic(err)
	}
	// 方法调用
	resp, err := client.RPCHello(context.Background(), &service.HelloRequest{
		MyName: "bob",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("message", resp.Message)
}
