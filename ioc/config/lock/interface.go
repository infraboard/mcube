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
	return ioc.Config().Get(LOCK).(*config).lf
}

type LockFactory interface {
	New(key string, ttl time.Duration) Lock
}

type Lock interface {
	// 锁配置
	WithOpt(opt *Options) Lock
	// 获取锁
	Lock(ctx context.Context) error
	// 释放锁
	UnLock(ctx context.Context) error
	// 刷新锁
	Refresh(ctx context.Context, ttl time.Duration) error
}
