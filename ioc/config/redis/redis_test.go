package redis_test

import (
	"testing"

	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/ioc/config/redis"
)

func TestRedisClient(t *testing.T) {
	m := redis.Client()
	t.Log(m)
}

func init() {

	err := ioc.ConfigIocObject(ioc.NewLoadConfigRequest())
	if err != nil {
		panic(err)
	}
}
