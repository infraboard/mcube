package ioc

const (
	DEFAULT_NAMESPACE = "default"
)

// 默认空间, 用于托管工具类, 在控制器之前进行初始化
func Default() StoreUser {
	return store.Namespace(DEFAULT_NAMESPACE)
}
