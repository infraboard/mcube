package cache

import (
	"github.com/infraboard/mcube/v2/ioc"

	"github.com/infraboard/mcube/v2/ioc/config/gocache"
	ioc_redis "github.com/infraboard/mcube/v2/ioc/config/redis"
)

func init() {
	ioc.Config().Registry(defaultConfig)
}

var defaultConfig = &cache{
	PROVIDER: PROVIDER_GO_CACHE,
	TTL:      300,
}

// Config 配置选项
type cache struct {
	// 使用换成提供方, 默认使用GoCache提供的内存缓存
	PROVIDER `json:"provider" yaml:"provider" toml:"provider" env:"PROVIDER"`
	// 单位秒, 默认5分钟
	TTL int64 `json:"ttl" yaml:"ttl" toml:"ttl" env:"TTL"`

	c Cache
	ioc.ObjectImpl
}

func (m *cache) Name() string {
	return AppName
}

func (m *cache) Priority() int {
	return 599
}

func (m *cache) Init() error {
	switch m.PROVIDER {
	case PROVIDER_REDIS:
		m.c = &redisCache{
			redis: ioc_redis.Client(),
			ttl:   m.TTL,
		}
	default:
		m.c = &goCache{
			gc:  gocache.C(),
			ttl: m.TTL,
		}
	}

	return nil
}
