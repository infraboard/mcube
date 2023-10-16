package cache

import (
	"context"

	"github.com/bluele/gcache"
	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	redis redis.UniversalClient
	ttl   int
}

func (r *redisCache) Set(
	ctx context.Context,
	key string,
	value any,
	opts ...SetOption) error {
	options := newOptions(r.ttl, opts...)
	return r.redis.Set(ctx, key, value, options.GetTTL()).Err()
}

func (r *redisCache) Del(ctx context.Context, keys ...string) error {
	return r.redis.Del(ctx, keys...).Err()
}

type goCache struct {
	gc  gcache.Cache
	ttl int
}

func (r *goCache) Set(
	ctx context.Context,
	key string,
	value any,
	opts ...SetOption) error {
	options := newOptions(r.ttl, opts...)
	return r.gc.SetWithExpire(key, value, options.GetTTL())
}

func (r *goCache) Del(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		r.gc.Remove(key)
	}
	return nil
}
