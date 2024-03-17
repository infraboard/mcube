package gocache

import (
	"github.com/bluele/gcache"
	"github.com/infraboard/mcube/v2/ioc"
)

const (
	AppName = "go_cache"
)

func C() gcache.Cache {
	return Get().c
}

func Get() *cache {
	obj := ioc.Config().Get(AppName)
	if obj == nil {
		return defaultConfig
	}
	return obj.(*cache)
}
