package lock

import (
	"github.com/infraboard/mcube/v2/ioc"
)

func init() {
	ioc.Config().Registry(defaultConfig)
}

var defaultConfig = &config{
	PROVIDER: PROVIDER_REDIS,
}

// Config 配置选项
type config struct {
	// 使用换成提供方, 默认使用REDIS提供分布式
	PROVIDER `json:"provider" yaml:"provider" toml:"provider" env:"LOCK_PROVIDER"`

	lf LockFactory

	ioc.ObjectImpl
}

func (c *config) Name() string {
	return AppName
}

func (c *config) Init() error {
	switch c.PROVIDER {
	case PROVIDER_REDIS:
		c.lf = NewRedisLockProvider()
	}
	return nil
}
