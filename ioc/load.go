package ioc

var (
	DefaultStore = &defaultStore{
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

func DevelopmentSetupWithPath(path string) {
	req := NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	req.ConfigFile.Path = path
	err := ConfigIocObject(req)
	if err != nil {
		panic(err)
	}
}

var (
	isLoaded bool
)

func ConfigIocObject(req *LoadConfigRequest) error {
	if isLoaded && !req.ForceLoad {
		return nil
	}

	// 加载对象的配置
	err := DefaultStore.LoadConfig(req)
	if err != nil {
		return err
	}

	// 初始化对象
	err = DefaultStore.InitIocObject()
	if err != nil {
		return err
	}

	// 依赖自动注入
	err = DefaultStore.Autowire()
	if err != nil {
		return err
	}

	isLoaded = true
	return nil
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
	// 默认加载后, 不允许重复加载, 这是为了避免多次初始化可能引发的问题
	ForceLoad bool
	// 环境变量配置
	ConfigEnv *configEnv
	// 文件配置方式
	ConfigFile *configFile
}

type configFile struct {
	Enabled bool
	Path    string
	// 如果找不到是否忽略
	SkipIFNotExist bool
}

type configEnv struct {
	Enabled bool
	Prefix  string
}
