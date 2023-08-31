package ioc

const (
	ConfigNamespace = "configs"
)

// 配置对象注册
func RegistryConfig(obj IocObject) {
	RegistryObjectWithNs(ConfigNamespace, obj)
}

// 获取配置对象
func GetConfig(name string) IocObject {
	return GetControllerWithVersion(name, DEFAULT_VERSION)
}

// 获取配置对象
func GetConfigWithVersion(name, version string) IocObject {
	return GetObjectWithNs(ConfigNamespace, name, version)
}
