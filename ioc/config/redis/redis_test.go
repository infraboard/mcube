package redis_test

import (
	"testing"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/redis"
	"github.com/infraboard/mcube/v2/tools/file"
)

func TestRedisClient(t *testing.T) {
	m := redis.Client()
	t.Log(m)
}

func TestDefaultConfig(t *testing.T) {
	file.MustToToml(
		redis.AppName,
		ioc.Config().Get(redis.AppName),
		"test/default.toml",
	)
}

func init() {

	err := ioc.ConfigIocObject(ioc.NewLoadConfigRequest())
	if err != nil {
		panic(err)
	}
}
