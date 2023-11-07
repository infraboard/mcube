package lock

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	luaRefresh = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("pexpire", KEYS[1], ARGV[2]) else return 0 end`)
	luaRelease = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("del", KEYS[1]) else return 0 end`)
	luaPTTL    = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("pttl", KEYS[1]) else return -3 end`)
	luaObtain  = redis.NewScript(`
if redis.call("set", KEYS[1], ARGV[1], "NX", "PX", ARGV[3]) then return redis.status_reply("OK") end

local offset = tonumber(ARGV[2])
if redis.call("getrange", KEYS[1], 0, offset-1) == string.sub(ARGV[1], 1, offset) then return redis.call("set", KEYS[1], ARGV[1], "PX", ARGV[3]) end
`)
)

type redisProvider struct {
	tmp   []byte
	tmpMu sync.Mutex
}

func (r *redisProvider) New(key string, ttl time.Duration, opts *Options) Lock {
	return &redisLock{
		client: nil,
	}
}

// Lock represents an obtained, distributed lock.
type redisLock struct {
	client   redis.Scripter
	key      string
	value    string
	tokenLen int
}

// 获取锁
func (l *redisLock) Lock(ctx context.Context) error {
	return nil
}

// 释放锁
func (l *redisLock) UnLock(ctx context.Context) error {
	return nil
}

// 刷新锁
func (l *redisLock) Refresh(ctx context.Context, ttl time.Duration) error {
	return nil
}
