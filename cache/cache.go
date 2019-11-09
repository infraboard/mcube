package cache

import (
	"fmt"
	"time"
)

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
	// Config the provider.
	Config(config string) error
	// close cache
	Close() error
}

// Instance is a function create a new Cache Instance
type Instance func() Cache

// NewCache Create a new cache driver by adapter name and config string.
// config need to be correct JSON as string: {"interval":360}.
func NewCache(adapterName, config string) (adapter Cache, err error) {
	instanceFunc, ok := adapters[adapterName]
	if !ok {
		err = fmt.Errorf("cache: unknown adapter name %q (forgot to import?)", adapterName)
		return
	}
	adapter = instanceFunc()
	err = adapter.Config(config)
	if err != nil {
		adapter = nil
	}
	return
}
