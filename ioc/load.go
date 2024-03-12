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

var (
	isLoaded bool
)

func ConfigIocObject(req *LoadConfigRequest) error {
	if isLoaded && !req.ForceLoad {
		return nil
	}

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
	err = store.Autowire()
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
}

type configEnv struct {
	Enabled bool
	Prefix  string
}
