package ioc

import (
	"google.golang.org/grpc"
)

const (
	controllerNamespace = "controllers"
)

// 用于托管控制器对象的Ioc空间, 配置完成后初始化
func Controller() StoreUser {
	return store.Namespace(controllerNamespace)
}

// GRPCService GRPC服务的实例
type GRPCControllerObject interface {
	Object
	Registry(*grpc.Server)
}

// 控制器对象注册
func RegistryController(obj Object) {
	RegistryObjectWithNs(controllerNamespace, obj)
}

// 获取控制器对象
func GetController(name string) Object {
	return GetObjectWithNs(controllerNamespace, name)
}

// LoadGrpcApp 加载所有的Grpc app
func LoadGrpcController(server *grpc.Server) {
	objects := store.Namespace(controllerNamespace)
	objects.ForEach(func(obj Object) {
		c, ok := obj.(GRPCControllerObject)
		if !ok {
			return
		}
		c.Registry(server)
	})
}
