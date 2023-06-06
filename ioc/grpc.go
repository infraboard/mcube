package ioc

import (
	"fmt"

	"google.golang.org/grpc"
)

var (
	grpcServers = map[string]GRPCServiceObject{}
)

// GRPCService GRPC服务的实例
type GRPCServiceObject interface {
	IocObject
	Registry(*grpc.Server)
}

// RegistryService 服务实例注册
func RegistryGrpcService(obj GRPCServiceObject) {
	// 已经注册的服务禁止再次注册
	_, ok := grpcServers[obj.Name()]
	if ok {
		panic(fmt.Sprintf("grpc app %s has registed", obj.Name()))
	}

	grpcServers[obj.Name()] = obj
}

// LoadedGrpcApp 查询加载成功的服务
func LoadedGrpcApp() (apps []string) {
	for k := range grpcServers {
		apps = append(apps, k)
	}
	return
}

func GetGrpcApp(name string) GRPCServiceObject {
	app, ok := grpcServers[name]
	if !ok {
		panic(fmt.Sprintf("grpc app %s not registed", name))
	}

	return app
}

// LoadGrpcApp 加载所有的Grpc app
func LoadGrpcApp(server *grpc.Server) {
	for _, app := range grpcServers {
		app.Registry(server)
	}
}
