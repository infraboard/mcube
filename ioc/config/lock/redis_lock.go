package lock

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"
	"strconv"
	"sync"
	"time"

	ioc_redis "github.com/infraboard/mcube/ioc/config/redis"
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

func NewRedisLockProvider() *RedisLockProvider {
	return &RedisLockProvider{}
}

type RedisLockProvider struct {
}

func (r *RedisLockProvider) New(key string, ttl time.Duration) Lock {
	return &RedisLock{
		client: ioc_redis.Client(),
		key:    key,
		ttl:    ttl,
		opt:    DefaultOptions(),
	}
}

// Lock represents an obtained, distributed lock.
type RedisLock struct {
	client   redis.Scripter
	key      string
	ttl      time.Duration
	opt      *Options
	value    string
	tokenLen int
	tmp      []byte
	tmpMu    sync.Mutex
}

// 锁配置
func (l *RedisLock) WithOpt(opt *Options) Lock {
	l.opt = opt
	return l
}

// 获取锁
func (l *RedisLock) Lock(ctx context.Context) error {
	token := l.opt.getToken()

	// Create a random token
	if token == "" {
		var err error
		if token, err = l.randomToken(); err != nil {
			return err
		}
	}

	value := token + l.opt.getMetadata()
	ttlVal := strconv.FormatInt(int64(l.ttl/time.Millisecond), 10)
	retry := l.opt.getRetryStrategy()

	// make sure we don't retry forever
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, time.Now().Add(l.ttl))
		defer cancel()
	}

	var ticker *time.Ticker
	for {
		ok, err := l.obtain(ctx, l.key, value, len(token), ttlVal)
		if err != nil {
			return err
		} else if ok {
			return nil
		}

		backoff := retry.NextBackoff()
		if backoff < 1 {
			return ErrNotObtained
		}

		if ticker == nil {
			ticker = time.NewTicker(backoff)
			defer ticker.Stop()
		} else {
			ticker.Reset(backoff)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func (c *RedisLock) obtain(ctx context.Context, key, value string, tokenLen int, ttlVal string) (bool, error) {
	_, err := luaObtain.Run(ctx, c.client, []string{key}, value, tokenLen, ttlVal).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (c *RedisLock) randomToken() (string, error) {
	c.tmpMu.Lock()
	defer c.tmpMu.Unlock()

	if len(c.tmp) == 0 {
		c.tmp = make([]byte, 16)
	}

	if _, err := io.ReadFull(rand.Reader, c.tmp); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(c.tmp), nil
}

// 释放锁
func (l *RedisLock) UnLock(ctx context.Context) error {
	if l == nil {
		return ErrLockNotHeld
	}

	res, err := luaRelease.Run(ctx, l.client, []string{l.key}, l.value).Result()
	if err == redis.Nil {
		return ErrLockNotHeld
	} else if err != nil {
		return err
	}

	if i, ok := res.(int64); !ok || i != 1 {
		return ErrLockNotHeld
	}
	return nil
}

// 刷新锁
func (l *RedisLock) Refresh(ctx context.Context, ttl time.Duration) error {
	ttlVal := strconv.FormatInt(int64(ttl/time.Millisecond), 10)
	status, err := luaRefresh.Run(ctx, l.client, []string{l.key}, l.value, ttlVal).Result()
	if err != nil {
		return err
	} else if status == int64(1) {
		return nil
	}
	return ErrNotObtained
}

// Key returns the redis key used by the lock.
func (l *RedisLock) Key() string {
	return l.key
}

// Token returns the token value set by the lock.
func (l *RedisLock) Token() string {
	return l.value[:l.tokenLen]
}

// Metadata returns the metadata of the lock.
func (l *RedisLock) Metadata() string {
	return l.value[l.tokenLen:]
}

// TTL returns the remaining time-to-live. Returns 0 if the lock has expired.
func (l *RedisLock) TTL(ctx context.Context) (time.Duration, error) {
	res, err := luaPTTL.Run(ctx, l.client, []string{l.key}, l.value).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	if num := res.(int64); num > 0 {
		return time.Duration(num) * time.Millisecond, nil
	}
	return 0, nil
}
