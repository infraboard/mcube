package redis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"

	"github.com/infraboard/mcube/cache"
)

// NewCache new an redis cache instance
func NewCache(conf *Config) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Address,
		Password: conf.Password,
		DB:       conf.DB,
	})

	return &Cache{
		prefix: conf.Prefix,
		ttl:    time.Duration(conf.DefaultTTL) * time.Second,
		client: client,
	}
}

// Cache 缓存
type Cache struct {
	prefix string
	ttl    time.Duration
	client *redis.Client
}

// SetDefaultTTL todo
func (c *Cache) SetDefaultTTL(ttl time.Duration) {
	c.ttl = ttl
}

// PutWithTTL todo
func (c *Cache) PutWithTTL(key string, val interface{}, ttl time.Duration) error {
	b, err := json.Marshal(val)
	if err != nil {
		return err
	}

	if err := c.client.Set(key, b, ttl).Err(); err != nil {
		return err
	}

	return nil
}

// Put todo
func (c *Cache) Put(key string, val interface{}) error {
	return c.PutWithTTL(key, val, c.ttl)
}

// Get todo
func (c *Cache) Get(key string, val interface{}) error {
	v, err := c.client.Get(key).Bytes()
	if err != nil {
		return err
	}

	if err := json.Unmarshal(v, val); err != nil {
		return err
	}

	return nil
}

// Delete todo
func (c *Cache) Delete(key string) error {
	if err := c.client.Del(key).Err(); err != nil {
		return err
	}

	return nil
}

// IsExist todo
func (c *Cache) IsExist(key string) bool {
	if ret := c.client.Exists(key).Val(); ret == 1 {
		return true
	}

	return false
}

// ClearAll todo
func (c *Cache) ClearAll() error {
	return c.client.FlushDB().Err()
}

// ListKey todo
func (c *Cache) ListKey(req *cache.ListKeyRequest) (*cache.ListKeyResponse, error) {
	fmt.Println(uint64(req.GetOffset()), req.Pattern(), int64(req.PageSize))
	ks, total, err := c.client.Scan(uint64(req.GetOffset()), req.Pattern(), int64(req.PageSize)).Result()
	if err != nil {
		return nil, err
	}
	return cache.NewListKeyResponse(ks, total), nil
}

// Decr todo
func (c *Cache) Decr(key string) error {
	return c.client.Decr(key).Err()
}

// Incr todo
func (c *Cache) Incr(key string) error {
	return c.client.Incr(key).Err()
}

// Close todo
func (c *Cache) Close() error {
	return c.client.Close()
}
