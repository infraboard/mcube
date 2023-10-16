package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/bluele/gcache"
	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	redis redis.UniversalClient
	ttl   int64
}

func (r *redisCache) Set(
	ctx context.Context,
	key string,
	value any,
	opts ...SetOption) error {
	options := newOptions(r.ttl, opts...)
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.redis.Set(ctx, key, data, options.GetTTL()).Err()
}

func (r *redisCache) Get(
	ctx context.Context,
	key string,
	value any) error {
	b := []byte{}
	res := r.redis.Get(ctx, key)
	if err := res.Err(); err != nil {
		return err
	}

	if err := res.Scan(b); err != nil {
		return err
	}
	return json.Unmarshal(b, value)
}

func (r *redisCache) Exist(
	ctx context.Context,
	key string) error {
	res := r.redis.Exists(ctx, key)
	if err := res.Err(); err != nil {
		return err
	}
	return nil
}

func (r *redisCache) Del(
	ctx context.Context,
	keys ...string) error {
	return r.redis.Del(ctx, keys...).Err()
}

func (r *redisCache) IncrBy(
	ctx context.Context,
	key string,
	value int64) (
	int64, error) {
	res := r.redis.IncrBy(ctx, key, value)
	if err := res.Err(); err != nil {
		return 0, err
	}
	return res.Val(), nil
}

type goCache struct {
	gc   gcache.Cache
	ttl  int64
	lock sync.Mutex
}

func (r *goCache) Set(
	ctx context.Context,
	key string,
	value any,
	opts ...SetOption) error {
	options := newOptions(r.ttl, opts...)
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.gc.SetWithExpire(key, b, options.GetTTL())
}

func (r *goCache) IncrBy(
	ctx context.Context,
	key string,
	value int64) (
	int64, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	var v int64
	err := r.Get(ctx, key, &v)
	if err != nil {
		if err == gcache.KeyNotFoundError {
			err := r.Set(ctx, key, value)
			if err != nil {
				return 0, err
			}
			return value, nil
		}
		return 0, err
	}

	v = v + value
	err = r.Set(ctx, key, v)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (r *goCache) Get(ctx context.Context, key string, value any) error {
	data, err := r.gc.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(data.([]byte), value)
}

func (r *goCache) Exist(ctx context.Context, key string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	if !r.gc.Has(key) {
		return fmt.Errorf("not exist")
	}
	return nil
}

func (r *goCache) Del(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		r.gc.Remove(key)
	}
	return nil
}
