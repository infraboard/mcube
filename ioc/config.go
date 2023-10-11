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

func NewLoadConfigRequest() *LoadConfigRequest {
	return &LoadConfigRequest{
		ConfigEnv: &ConfigEnv{
			Enabled: true,
		},
		ConfigFile: &ConfigFile{
			Enabled: false,
			Path:    "etc/application.toml",
		},
	}
}

type LoadConfigRequest struct {
	// 环境变量配置
	ConfigEnv *ConfigEnv
	// 文件配置方式
	ConfigFile *ConfigFile
}

type ConfigFile struct {
	Enabled bool
	Path    string
}

type ConfigEnv struct {
	Enabled bool
	Prefix  string
}
