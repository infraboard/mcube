package ioc

const (
	configNamespace = "configs"
)

// 用于托管配置对象的Ioc空间, 最先初始化
func Config() StoreUser {
	return store.
		Namespace(configNamespace).
		SetPriority(99)
}
