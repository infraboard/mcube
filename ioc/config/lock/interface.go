package lock

import (
	"context"
	"errors"
	"time"

	"github.com/infraboard/mcube/ioc"
)

const (
	LOCK = "lock"
)

var (
	// ErrNotObtained is returned when a lock cannot be obtained.
	ErrNotObtained = errors.New("lock: not obtained")

	// ErrLockNotHeld is returned when trying to release an inactive lock.
	ErrLockNotHeld = errors.New("lock: lock not held")
)

func L() LockFactory {
	return ioc.Config().Get(LOCK).(*config).LockFactory
}

type LockFactory interface {
	New(key string, ttl time.Duration, opts *Options) Lock
}

type Lock interface {
	// 获取锁
	Lock(ctx context.Context) error
	// 释放锁
	UnLock(ctx context.Context) error
	// 刷新锁
	Refresh(ctx context.Context, ttl time.Duration) error
}

// Options describe the options for the lock
type Options struct {
	// RetryStrategy allows to customise the lock retry strategy.
	// Default: do not retry
	RetryStrategy RetryStrategy

	// Metadata string.
	Metadata string

	// Token is a unique value that is used to identify the lock. By default, a random tokens are generated. Use this
	// option to provide a custom token instead.
	Token string
}

func (o *Options) getMetadata() string {
	if o != nil {
		return o.Metadata
	}
	return ""
}

func (o *Options) getToken() string {
	if o != nil {
		return o.Token
	}
	return ""
}

func (o *Options) getRetryStrategy() RetryStrategy {
	if o != nil && o.RetryStrategy != nil {
		return o.RetryStrategy
	}
	return NoRetry()
}
