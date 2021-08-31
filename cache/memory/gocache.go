package memory

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/bluele/gcache"

	"github.com/infraboard/mcube/cache"
)

// NewCache new an redis cache instance
func NewCache(conf *Config) *Cache {
	gc := gcache.New(conf.Size).
		LRU().
		Expiration(conf.GetTTL()).
		Build()

	return &Cache{
		ttl:  conf.GetTTL(),
		size: conf.Size,
		gc:   gc,
	}
}

// Cache 内存缓存
type Cache struct {
	ttl  time.Duration
	size int
	gc   gcache.Cache
	ctx  context.Context
}

// SetDefaultTTL 修复默认TTL
func (c *Cache) SetDefaultTTL(ttl time.Duration) {
	c.ttl = ttl
}

// Put 添加key
func (c *Cache) Put(key string, val interface{}) error {
	return c.PutWithTTL(key, val, c.ttl)
}

// PutWithTTL 指定TTL
func (c *Cache) PutWithTTL(key string, val interface{}, ttl time.Duration) error {
	b, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return c.gc.SetWithExpire(key, b, c.ttl)
}

// Get 获取对象
func (c *Cache) Get(key string, val interface{}) error {
	v, err := c.gc.Get(key)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(v.([]byte), val); err != nil {
		return err
	}

	return nil
}

// Delete 删除对象
func (c *Cache) Delete(key string) error {
	if !c.gc.Remove(key) {
		return errors.New("delete cache error")
	}

	return nil
}

// IsExist 判断Key是否存在
func (c *Cache) IsExist(key string) bool {
	return c.gc.Has(key)
}

// ClearAll 清除缓存
func (c *Cache) ClearAll() error {
	c.gc.Purge()
	return nil
}

// Incr 等待实现
func (c *Cache) Incr(key string) error {
	return errors.New("not impliment")
}

// Decr 等待实现
func (c *Cache) Decr(key string) error {
	return errors.New("not impliment")
}

// Close 关闭缓存
func (c *Cache) Close() error {
	return nil
}

// ListKey todo
func (c *Cache) ListKey(req *cache.ListKeyRequest) (*cache.ListKeyResponse, error) {
	ks := []string{}
	for _, k := range c.gc.Keys(true) {
		if v, ok := k.(string); ok {
			ks = append(ks, v)
		}
	}

	return cache.NewListKeyResponse(ks, uint64(len(ks))), nil
}

// WithContext 携带上下文
func (c *Cache) WithContext(ctx context.Context) cache.Cache {
	clone := *c
	clone.ctx = ctx
	return &clone
}
