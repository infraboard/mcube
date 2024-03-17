package redis

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/redis/go-redis/v9"
)

const (
	AppName = "redis"
)

func Client() redis.UniversalClient {
	return Get().client
}

func Get() *Redist {
	obj := ioc.Config().Get(AppName)
	if obj == nil {
		return defaultConfig
	}
	return obj.(*Redist)
}
