package flowcontrol

import (
	"context"
)

// RateLimiter 限流器
type RateLimiter interface {
	// TryAccept returns true if a token is taken immediately. Otherwise,
	// it returns false.
	TryAccept() bool
	// Accept returns once a token becomes available.
	Accept()
	// Stop stops the rate limiter, subsequent calls to CanAccept will return false
	Stop()
	// QPS returns QPS of this rate limiter
	QPS() float32
	// Wait returns nil if a token is taken before the Context is done.
	Wait(ctx context.Context) error
}
