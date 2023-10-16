package cache

import (
	"github.com/infraboard/mcube/ioc"

	"github.com/infraboard/mcube/ioc/config/gocache"
	ioc_redis "github.com/infraboard/mcube/ioc/config/redis"
)

func init() {
	ioc.Config().Registry(&cache{
		PROVIDER: PROVIDER_GO_CACHE,
		TTL:      300,
	})
}

// Config 配置选项
type cache struct {
	// 使用换成提供方, 默认使用GoCache提供的内存缓存
	PROVIDER `json:"provider" yaml:"provider" toml:"provider" env:"CACHE_PROVIDER"`
	// 单位秒, 默认5分钟
	TTL int64 `json:"ttl" yaml:"ttl" toml:"ttl" env:"CACHE_TTL"`

	c Cache
	ioc.ObjectImpl
}

func (m *cache) Name() string {
	return CACHE
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
