package redis

import (
	"github.com/infraboard/mcube/ioc"
	"github.com/redis/go-redis/v9"
)

const (
	REDIS = "redis"
)

func Client() redis.UniversalClient {
	return ioc.Config().Get(REDIS).(*Redist).client
}
