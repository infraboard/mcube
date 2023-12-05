package redis

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/redis/go-redis/v9"
)

const (
	AppName = "redis"
)

func Client() redis.UniversalClient {
	return ioc.Config().Get(AppName).(*Redist).client
}
