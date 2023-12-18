package main

import (
	"context"
	"fmt"

	"github.com/infraboard/mcube/v2/examples/rpc_grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 连接到服务
	conn, err := grpc.Dial(
		"127.0.0.1:18080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	resp, err := pb.NewHelloClient(conn).Greet(context.Background(), &pb.GreetRequest{
		Name: "old yu",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
