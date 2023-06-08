package cache

import (
	"time"

	"github.com/infraboard/mcube/http/request"
)

var (
	cache Cache
)

// C 全局缓存对象, 默认使用
func C() Cache {
	if cache == nil {
		panic("global cache instance is nil")
	}
	return cache
}

// SetGlobal 设置全局缓存
func SetGlobal(c Cache) {
	cache = c
}

// Cache provides the interface for cache implementations.
type Cache interface {
	ListKey(*ListKeyRequest) (*ListKeyResponse, error)
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

// NewListKeyRequest todo
func NewListKeyRequest(pattern string, ps uint, pn uint) *ListKeyRequest {
	return &ListKeyRequest{
		pattern:     pattern,
		PageRequest: request.NewPageRequest(ps, pn),
	}
}

// ListKeyRequest todo
type ListKeyRequest struct {
	pattern string
	*request.PageRequest
}

// Pattern tood
func (req *ListKeyRequest) Pattern() string {
	return req.pattern
}

// NewListKeyResponse todo
func NewListKeyResponse(keys []string, total uint64) *ListKeyResponse {
	return &ListKeyResponse{
		Keys:  keys,
		Total: total,
	}
}

// ListKeyResponse todo
type ListKeyResponse struct {
	Keys  []string
	Total uint64
}
