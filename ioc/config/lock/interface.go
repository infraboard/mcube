package lock

import (
	"context"
	"errors"
	"time"

	"github.com/infraboard/mcube/v2/ioc"
)

const (
	AppName = "lock"
)

var (
	// ErrNotObtained is returned when a lock cannot be obtained.
	ErrNotObtained = errors.New("lock: not obtained")

	// ErrLockNotHeld is returned when trying to release an inactive lock.
	ErrLockNotHeld = errors.New("lock: lock not held")
)

func L() LockFactory {
	return Get().lf
}

func Get() *config {
	obj := ioc.Config().Get(AppName)
	if obj == nil {
		return defaultConfig
	}
	return obj.(*config)
}

type LockFactory interface {
	New(key string, ttl time.Duration) Lock
}

type Lock interface {
	// 锁配置
	WithOpt(opt *Options) Lock
	// TryLock
	TryLock(ctx context.Context) error
	// 获取锁
	Lock(ctx context.Context) error
	// 释放锁
	UnLock(ctx context.Context) error
	// 刷新锁
	Refresh(ctx context.Context, ttl time.Duration) error
}
