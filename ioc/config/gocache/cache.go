package gocache

import (
	"github.com/bluele/gcache"
	"github.com/infraboard/mcube/v2/ioc"
)

func init() {
	ioc.Config().Registry(defaultConfig)
}

var defaultConfig = &cache{
	Size: 500,
}

// Config 配置选项
type cache struct {
	// 个数
	Size int `json:"size" yaml:"size" toml:"size" env:"GO_CACHE_SIZE"`

	c gcache.Cache
	ioc.ObjectImpl
}

func (m *cache) Name() string {
	return AppName
}

func (i *cache) Priority() int {
	return 696
}

func (m *cache) Init() error {
	m.c = gcache.New(m.Size).
		LRU().
		Build()
	return nil
}
