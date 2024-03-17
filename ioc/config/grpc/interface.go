package grpc

import "github.com/infraboard/mcube/v2/ioc"

const (
	AppName = "grpc"
)

func Get() *Grpc {
	obj := ioc.Config().Get(AppName)
	if obj == nil {
		return defaultConfig
	}
	return obj.(*Grpc)
}
