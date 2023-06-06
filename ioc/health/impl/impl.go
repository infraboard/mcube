package impl

import (
	"github.com/infraboard/mcube/ioc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"

	my "github.com/infraboard/mcube/ioc/health"
)

var (
	// Service 服务实例
	svr = &impl{Server: health.NewServer()}
)

type impl struct {
	*health.Server
	ioc.IocObjectImpl
}

func (i *impl) Init() error {
	return nil
}

func (i *impl) Name() string {
	return my.AppName
}

func (i *impl) Registry(server *grpc.Server) {
	healthgrpc.RegisterHealthServer(server, svr)
}

func init() {
	ioc.RegistryGrpcService(svr)
}
