package cache

import (
	"time"

	"github.com/infraboard/mcube/cache/memory"
)

var (
	cache Cache
)

// C 全局缓存对象, 默认使用
func C() Cache {
	if cache == nil {
		mc := memory.NewCache(memory.NewDefaultConfig())
		SetGlobal(mc)
	}
	return cache
}

// SetGlobal 设置全局缓存
func SetGlobal(c Cache) {
	cache = c
}

// Cache provides the interface for cache implementations.
type Cache interface {
	// SetTTL set default ttl
	SetDefaultTTL(ttl time.Duration)
	// set cached value with key and expire time.
	Put(key string, val interface{}) error
	// set cached value with key and expire time.
	PutWithTTL(key string, val interface{}, ttl time.Duration) error
	// get cached value by key.
	Get(key string, val interface{}) error
	// delete cached value by key.
	Delete(key string) error
	// check if cached value exists or not.
	IsExist(key string) bool
	// clear all cache.
	ClearAll() error
	// increase cached int value by key, as a counter.
	Incr(key string) error
	// decrease cached int value by key, as a counter.
	Decr(key string) error
	// close cache
	Close() error
}
