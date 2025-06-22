package service

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func NewClient(address string) (HelloService, error) {
	// 建立TCP连接
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("net.Dial:", err)
	}

	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))

	// buf := bytes.NewBuffer([]byte{})
	// gob.NewEncoder(buf).Encode(HelloRequest{Name: "test"})

	//
	// reader := bytes.NewReader([]byte("{}"))
	// resp := HelloRequest{Name: "test"}
	// gob.NewDecoder(reader).Decode(&resp)
	return &HelloServiceClient{conn: client}, nil
}

// 要封装原始的 不友好的rpc call
type HelloServiceClient struct {
	conn *rpc.Client
}

func (c *HelloServiceClient) Hello(request *HelloRequest, response *HelloResponse) error {
	if err := c.conn.Call("HelloService.Hello", request, response); err != nil {
		panic(err)
	}
	return nil
}
