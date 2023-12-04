package gocache

import (
	"github.com/bluele/gcache"
	"github.com/infraboard/mcube/ioc"
)

const (
	AppName = "go_cache"
)

func C() gcache.Cache {
	return ioc.Config().Get(AppName).(*cache).c
}
