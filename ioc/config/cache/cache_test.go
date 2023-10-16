package cache_test

import (
	"context"
	"testing"

	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/ioc/config/cache"
)

var (
	ctx = context.Background()
)

func TestGetClientGetter(t *testing.T) {
	m := cache.C()
	t.Log(m.Set(ctx, "a", "v1"))

	var a string
	t.Log(m.Get(ctx, "a", &a))
	t.Log(a)

	t.Log(m.IncrBy(ctx, "b", 1))
	t.Log(m.IncrBy(ctx, "b", 1))
	t.Log(m.IncrBy(ctx, "b", 1))
	var b int64
	t.Log(m.Get(ctx, "b", &b))
	t.Log(b)
}

func init() {
	err := ioc.ConfigIocObject(ioc.NewLoadConfigRequest())
	if err != nil {
		panic(err)
	}
}
