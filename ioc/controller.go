package ioc

import (
	"google.golang.org/grpc"
)

const (
	CONTROLLER_NAMESPACE = "controllers"
)

// 用于托管控制器对象的Ioc空间, 配置完成后初始化
func Controller() StoreUser {
	return store.Namespace(CONTROLLER_NAMESPACE)
}

// GRPCService GRPC服务的实例
type GRPCControllerObject interface {
	Object
	Registry(*grpc.Server)
}

// 控制器对象注册
func RegistryController(obj Object) {
	RegistryObjectWithNs(CONTROLLER_NAMESPACE, obj)
}

// 获取控制器对象
func GetController(name string) Object {
	return GetObjectWithNs(CONTROLLER_NAMESPACE, name)
}

// LoadGrpcApp 加载所有的Grpc app
func LoadGrpcController(server *grpc.Server) {
	objects := store.Namespace(CONTROLLER_NAMESPACE)
	objects.ForEach(func(w *ObjectWrapper) {
		c, ok := w.Value.(GRPCControllerObject)
		if !ok {
			return
		}
		c.Registry(server)
	})
}

func GrpcControllerCount() int {
	count := 0
	objects := store.Namespace(CONTROLLER_NAMESPACE)
	objects.ForEach(func(w *ObjectWrapper) {
		_, ok := w.Value.(GRPCControllerObject)
		if ok {
			count++
		}
	})
	return count
}
