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
	// 使用默认配置文件 etc/application.toml
	req.ConfigFile.Enabled = true
	req.ConfigFile.Paths = []string{"etc/application.toml"}
	req.ConfigFile.SkipIFNotExist = true
	err := ConfigIocObject(req)
	if err != nil {
		panic(err)
	}
}

func DevelopmentSetupWithPath(path string) {
	req := NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	req.ConfigFile.Paths = []string{path}
	err := ConfigIocObject(req)
	if err != nil {
		panic(err)
	}
}

// DevelopmentSetupWithPaths 支持多个配置文件的开发环境初始化
func DevelopmentSetupWithPaths(paths ...string) {
	req := NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	req.ConfigFile.Paths = paths
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

	// 1. 加载对象的配置
	err := DefaultStore.LoadConfig(req)
	if err != nil {
		return err
	}

	// 2. 调用 PostConfig 钩子（配置验证）
	err = DefaultStore.CallPostConfigHooks()
	if err != nil {
		return err
	}

	// 3. 依赖自动注入
	err = DefaultStore.Autowire()
	if err != nil {
		return err
	}

	// 4. 初始化对象（包含 PreInit 和 PostInit 钩子）
	err = DefaultStore.InitIocObject()
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
			Paths:   []string{}, // 默认为空，需要显式指定配置文件
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
	// 配置文件路径列表，支持多个文件，后面的会覆盖前面的
	Paths []string
	// 如果找不到是否忽略
	SkipIFNotExist bool
}

// Path 获取第一个配置文件路径（向后兼容）
// Deprecated: 使用 Paths 代替
func (c *configFile) Path() string {
	if len(c.Paths) > 0 {
		return c.Paths[0]
	}
	return ""
}

// SetPath 设置单个配置文件路径（向后兼容）
// Deprecated: 使用 Paths 代替
func (c *configFile) SetPath(path string) {
	c.Paths = []string{path}
}

type configEnv struct {
	Enabled bool
	Prefix  string
}
