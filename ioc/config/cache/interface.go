package cache

import "context"

type Cache interface {
	Set(ctx context.Context, key string, value any, options ...SetOption)
	Del(ctx context.Context, keys ...string)
}
