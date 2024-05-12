package lock

import (
	"context"
	"time"

	"github.com/bluele/gcache"
	"github.com/infraboard/mcube/v2/ioc/config/gocache"
)

func NewGoCacheLockProvider() *GoCacheLockProvider {
	return &GoCacheLockProvider{}
}

type GoCacheLockProvider struct {
}

func (r *GoCacheLockProvider) New(key string, ttl time.Duration) Lock {

	return &GoCacheLock{
		key:   key,
		ttl:   ttl,
		cache: gocache.C(),
		opt:   DefaultOptions(),
	}
}

type GoCacheLock struct {
	cache gcache.Cache
	key   string
	value string
	ttl   time.Duration
	opt   *Options
}

// 锁配置
func (m *GoCacheLock) WithOpt(opt *Options) Lock {
	return m
}

// 获取锁
func (m *GoCacheLock) Lock(ctx context.Context) error {
	retry := m.opt.getRetryStrategy()
	// make sure we don't retry forever
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, time.Now().Add(m.ttl))
		defer cancel()
	}

	var ticker *time.Ticker
	for {
		ok, err := m.obtain(ctx)
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

func (m *GoCacheLock) obtain(context.Context) (bool, error) {
	if m.cache.Has(m.key) {
		return false, nil
	}

	if err := m.cache.SetWithExpire(m.key, m.value, m.ttl); err != nil {
		return false, err
	}
	return true, nil
}

// 释放锁
func (m *GoCacheLock) UnLock(context.Context) error {
	m.cache.Remove(m.key)
	return nil
}

// 刷新锁
func (m *GoCacheLock) Refresh(ctx context.Context, ttl time.Duration) error {
	return m.cache.SetWithExpire(m.key, m.value, m.ttl)
}
