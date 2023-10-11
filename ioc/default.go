package ioc

const (
	defaultNamespace = "default"
)

// 默认空间, 用于托管工具类, 在控制器之前进行初始化
func Default() StoreUser {
	return store.
		Namespace(defaultNamespace).
		SetPriority(9)
}
