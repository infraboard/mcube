package service

import "context"

const (
	APP_NAME = "HelloService"
)

type HelloService interface {
	RPCHello(context.Context, *HelloRequest) (*HelloResponse, error)
}

type HelloRequest struct {
	MyName string `json:"my_name"`
}

type HelloResponse struct {
	Message string `json:"message"`
}
