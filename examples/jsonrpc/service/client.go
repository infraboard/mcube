package service

import (
	"fmt"

	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	"github.com/infraboard/mcube/v2/ioc/config/jsonrpc"
	"resty.dev/v3"
)

func NewClient(address string) (HelloService, error) {
	// 建立TCP连接
	client := resty.New().
		SetDebug(application.Get().Debug).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBaseURL(address).AddResponseMiddleware(func(c *resty.Client, r *resty.Response) error {
		if r.StatusCode()/100 != 2 {
			return exception.NewApiExceptionFromString(r.String())
		}
		return nil
	})
	return &HelloServiceClient{client: client}, nil
}

// 要封装原始的 不友好的rpc call
type HelloServiceClient struct {
	client *resty.Client
}

func (c *HelloServiceClient) Hello(in *HelloRequest, out *HelloResponse) error {
	body := jsonrpc.NewRequest(fmt.Sprintf("%s.Hello", APP_NAME), in)
	result := jsonrpc.NewResponse(body.Id, out)
	_, err := c.client.R().SetDebug(true).SetBody(body).SetResult(result).Post("")
	if err != nil {
		return err
	}

	if result.Error != nil {
		return exception.NewApiExceptionFromString(*result.Error)
	}
	return nil
}
