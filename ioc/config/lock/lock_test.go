package lock_test

import (
	"context"
	"testing"
	"time"

	"github.com/infraboard/mcube/ioc/config/lock"
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
