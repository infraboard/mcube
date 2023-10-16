package gocache

import (
	"github.com/bluele/gcache"
	"github.com/infraboard/mcube/ioc"
)

const (
	GO_CACHE = "go_cache"
)

func C() gcache.Cache {
	return ioc.Config().Get(GO_CACHE).(*cache).c
}
