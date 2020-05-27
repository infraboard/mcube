package flowcontrol

import "time"

// RateLimiter 限流器
type RateLimiter interface {
	Wait(count int64)
	WaitMaxDuration(count int64, maxWait time.Duration) bool
	Take(count int64) time.Duration
	TakeMaxDuration(count int64, maxWait time.Duration) (time.Duration, bool)
	TakeAvailable(count int64) int64
	TakeOneAvailable() bool
	LastTakeTime() time.Time
	Available() int64
	Capacity() int64
	Rate() float64
}
