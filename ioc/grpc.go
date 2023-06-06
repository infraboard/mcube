package ioc

import (
	"google.golang.org/grpc"
)

var (
	GrpcServiceNamespace = "grpc"
)

// GRPCService GRPC服务的实例
type GRPCServiceObject interface {
	IocObject
	Registry(*grpc.Server)
}

// RegistryService 服务实例注册
func RegistryGrpcService(obj GRPCServiceObject) {
	RegistryObjectWithNs(GrpcServiceNamespace, obj)
}

// 获取gprc对象
func GetGrpcService(name string) IocObject {
	return GetObjectWithNs(GrpcServiceNamespace, name)
}

// LoadedGrpcService 查询加载成功的服务
func LoadedGrpcService() (names []string) {
	return store.Namespace(GrpcServiceNamespace).ObjectNames()
}

// LoadGrpcApp 加载所有的Grpc app
func LoadGrpcService(server *grpc.Server) {
	objects := store.Namespace(GrpcServiceNamespace)
	for _, obj := range objects.Items {
		obj.(GRPCServiceObject).Registry(server)
	}
}
