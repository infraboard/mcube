package ioc

const (
	CONTROLLER_NAMESPACE = "controllers"
)

// 用于托管控制器对象的Ioc空间, 配置完成后初始化
func Controller() StoreUser {
	return store.Namespace(CONTROLLER_NAMESPACE)
}
