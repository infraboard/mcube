package ioc

func DevelopmentSetup() {
	req := NewLoadConfigRequest()
	err := ConfigIocObject(req)
	if err != nil {
		panic(err)
	}
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
