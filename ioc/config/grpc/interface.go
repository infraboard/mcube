package grpc

import "github.com/infraboard/mcube/v2/ioc"

const (
	AppName = "grpc"
)

func Get() *Grpc {
	return ioc.Config().Get(AppName).(*Grpc)
}
