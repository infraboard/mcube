package cache

import (
	"github.com/infraboard/mcube/ioc/config/gocache"
	"github.com/infraboard/mcube/ioc/config/redis"
)

type PROVIDER string

const (
	PROVIDER_GO_CACHE = gocache.GO_CACHE
	PROVIDER_REDIS    = redis.REDIS
)
