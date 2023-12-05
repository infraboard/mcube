package gocache_test

import (
	"testing"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/gocache"
)

func TestGetClientGetter(t *testing.T) {
	m := gocache.C()
	t.Log(m)
}

func init() {
	err := ioc.ConfigIocObject(ioc.NewLoadConfigRequest())
	if err != nil {
		panic(err)
	}
}
