package ioc

import (
	"google.golang.org/grpc"
)

var (
	ControllerNamespace = "controllers"
)

// GRPCService GRPC服务的实例
type GRPCControllerObject interface {
	Registry(*grpc.Server)
}

// 控制器对象注册
func RegistryController(obj IocObject) {
	RegistryObjectWithNs(ControllerNamespace, obj)
}

// 获取控制器对象
func GetController(name string) IocObject {
	return GetObjectWithNs(ControllerNamespace, name)
}

// 或者注册完成的控制器
func ListControllerObjectNames() (names []string) {
	return store.Namespace(ControllerNamespace).ObjectNames()
}

// LoadGrpcApp 加载所有的Grpc app
func LoadGrpcController(server *grpc.Server) {
	objects := store.Namespace(ControllerNamespace)
	for _, obj := range objects.Items {
		c, ok := obj.(GRPCControllerObject)
		if !ok {
			continue
		}
		c.Registry(server)
	}
}
