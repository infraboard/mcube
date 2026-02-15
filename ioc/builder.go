package ioc

import (
	"path/filepath"
)

// ConfigLoader 配置加载器，提供流畅的Builder API
type ConfigLoader struct {
	req *LoadConfigRequest
}

// LoadConfig 创建配置加载器
// 使用示例:
//
//	err := ioc.LoadConfig().
//	    FromFile("etc/base.toml").
//	    FromFile("etc/prod.toml").
//	    FromEnv("APP").
//	    Load()
func LoadConfig() *ConfigLoader {
	return &ConfigLoader{
		req: NewLoadConfigRequest(),
	}
}

// FromFile 添加单个配置文件
func (c *ConfigLoader) FromFile(path string) *ConfigLoader {
	c.req.ConfigFile.Enabled = true
	c.req.ConfigFile.Paths = append(c.req.ConfigFile.Paths, path)
	return c
}

// FromFiles 添加多个配置文件
func (c *ConfigLoader) FromFiles(paths ...string) *ConfigLoader {
	c.req.ConfigFile.Enabled = true
	c.req.ConfigFile.Paths = append(c.req.ConfigFile.Paths, paths...)
	return c
}

// FromPattern 从glob模式添加配置文件
// 例如: FromPattern("etc/*.toml") 会加载etc目录下所有toml文件
func (c *ConfigLoader) FromPattern(pattern string) *ConfigLoader {
	matches, err := filepath.Glob(pattern)
	if err == nil && len(matches) > 0 {
		c.req.ConfigFile.Enabled = true
		c.req.ConfigFile.Paths = append(c.req.ConfigFile.Paths, matches...)
	}
	return c
}

// FromEnv 从环境变量加载配置
func (c *ConfigLoader) FromEnv(prefix string) *ConfigLoader {
	c.req.ConfigEnv.Enabled = true
	c.req.ConfigEnv.Prefix = prefix
	return c
}

// SkipIfNotExist 如果配置文件不存在则跳过，不报错
func (c *ConfigLoader) SkipIfNotExist() *ConfigLoader {
	c.req.ConfigFile.SkipIFNotExist = true
	return c
}

// ForceReload 强制重新加载，即使已经加载过
func (c *ConfigLoader) ForceReload() *ConfigLoader {
	c.req.ForceLoad = true
	return c
}

// Load 执行配置加载
func (c *ConfigLoader) Load() error {
	return ConfigIocObject(c.req)
}

// WithConfigFile 函数式选项：设置单个配置文件
func WithConfigFile(path string) func(*LoadConfigRequest) {
	return func(req *LoadConfigRequest) {
		req.ConfigFile.Enabled = true
		req.ConfigFile.Paths = []string{path}
	}
}

// WithConfigFiles 函数式选项：设置多个配置文件
func WithConfigFiles(paths ...string) func(*LoadConfigRequest) {
	return func(req *LoadConfigRequest) {
		req.ConfigFile.Enabled = true
		req.ConfigFile.Paths = paths
	}
}

// WithEnvPrefix 函数式选项：设置环境变量前缀
func WithEnvPrefix(prefix string) func(*LoadConfigRequest) {
	return func(req *LoadConfigRequest) {
		req.ConfigEnv.Enabled = true
		req.ConfigEnv.Prefix = prefix
	}
}

// SkipNotExist 函数式选项：跳过不存在的文件
func SkipNotExist() func(*LoadConfigRequest) {
	return func(req *LoadConfigRequest) {
		req.ConfigFile.SkipIFNotExist = true
	}
}

// Load 函数式风格的配置加载
// 使用示例:
//
//	err := ioc.Load(
//	    ioc.WithConfigFiles("etc/base.toml", "etc/prod.toml"),
//	    ioc.WithEnvPrefix("APP"),
//	    ioc.SkipNotExist(),
//	)
func Load(opts ...func(*LoadConfigRequest)) error {
	req := NewLoadConfigRequest()
	for _, opt := range opts {
		opt(req)
	}
	return ConfigIocObject(req)
}

// Setup 快速配置加载（推荐用于main函数）
// 使用示例:
//
//	// 开发环境
//	ioc.Setup(ioc.WithConfigFiles("etc/base.toml", "etc/dev.toml"))
//
//	// 生产环境
//	ioc.Setup(
//	    ioc.WithConfigFiles("etc/base.toml", "etc/prod.toml"),
//	    ioc.WithEnvPrefix("APP"),
//	)
func Setup(opts ...func(*LoadConfigRequest)) error {
	return Load(opts...)
}
