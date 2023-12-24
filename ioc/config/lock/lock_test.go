package lock_test

import (
	"context"
	"testing"
	"time"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/lock"
	"github.com/infraboard/mcube/v2/tools/file"
)

var (
	ctx = context.Background()
)

func TestRedisLock(t *testing.T) {
	m := lock.L().New("test", 10*time.Second)
	t.Log(m)

	if err := m.Lock(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestDefaultConfig(t *testing.T) {
	file.MustToToml(
		lock.AppName,
		ioc.Config().Get(lock.AppName),
		"test/default.toml",
	)
}
