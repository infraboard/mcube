package ioc

var (
	store = &defaultStore{
		store: []*NamespaceStore{
			newNamespaceStore(CONFIG_NAMESPACE).SetPriority(99),
			newNamespaceStore(CONTROLLER_NAMESPACE).SetPriority(0),
			newNamespaceStore(DEFAULT_NAMESPACE).SetPriority(9),
			newNamespaceStore(API_NAMESPACE).SetPriority(-99),
		},
	}
)

func DevelopmentSetup() {
	req := NewLoadConfigRequest()
	err := ConfigIocObject(req)
	if err != nil {
		panic(err)
	}
}

func ConfigIocObject(req *LoadConfigRequest) error {
	// 加载对象的配置
	err := store.LoadConfig(req)
	if err != nil {
		return err
	}

	// 初始化对象
	err = store.InitIocObject()
	if err != nil {
		return err
	}

	// 依赖自动注入
	return store.Autowire()
}

func NewLoadConfigRequest() *LoadConfigRequest {
	return &LoadConfigRequest{
		ConfigEnv: &configEnv{
			Enabled: true,
		},
		ConfigFile: &configFile{
			Enabled: false,
			Path:    "etc/application.toml",
		},
	}
}

type LoadConfigRequest struct {
	// 环境变量配置
	ConfigEnv *configEnv
	// 文件配置方式
	ConfigFile *configFile
}

type configFile struct {
	Enabled bool
	Path    string
}

type configEnv struct {
	Enabled bool
	Prefix  string
}
