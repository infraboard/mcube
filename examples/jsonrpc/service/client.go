package service

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func NewClient(address string) (HelloService, error) {
	// 建立TCP连接
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal("net.Dial:", err)
	}

	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
	return &HelloServiceClient{conn: client}, nil
}

// 要封装原始的 不友好的rpc call
type HelloServiceClient struct {
	conn *rpc.Client
}

func (c *HelloServiceClient) Hello(req *HelloRequest, resp *HelloResponse) error {
	fmt.Print(req, resp)
	if err := c.conn.Call("HelloService.Hello", req, resp); err != nil {
		return err
	}
	return nil
}
