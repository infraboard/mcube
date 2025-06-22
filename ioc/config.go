package ioc

const (
	CONFIG_NAMESPACE = "configs"
)

// 用于托管配置对象的Ioc空间, 最先初始化
func Config() StoreUser {
	return DefaultStore.Namespace(CONFIG_NAMESPACE)
}
