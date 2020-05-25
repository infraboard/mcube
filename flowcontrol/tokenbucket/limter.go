package tokenbucket

import (
	"context"
	"errors"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/infraboard/mcube/flowcontrol"
)

type tokenBucketRateLimiter struct {
	limiter *rate.Limiter
	clock   Clock
	qps     float32
}

// NewTokenBucketRateLimiter creates a rate limiter which implements a token bucket approach.
// The rate limiter allows bursts of up to 'burst' to exceed the QPS, while still maintaining a
// smoothed qps rate of 'qps'.
// The bucket is initially filled with 'burst' tokens, and refills at a rate of 'qps'.
// The maximum number of tokens in the bucket is capped at 'burst'.
func NewTokenBucketRateLimiter(qps float32, burst int) flowcontrol.RateLimiter {
	limiter := rate.NewLimiter(rate.Limit(qps), burst)
	return newTokenBucketRateLimiter(limiter, realClock{}, qps)
}

// Clock An injectable, mockable clock interface.
type Clock interface {
	Now() time.Time
	Sleep(time.Duration)
}

type realClock struct{}

func (realClock) Now() time.Time {
	return time.Now()
}
func (realClock) Sleep(d time.Duration) {
	time.Sleep(d)
}

// NewTokenBucketRateLimiterWithClock is identical to NewTokenBucketRateLimiter
// but allows an injectable clock, for testing.
func NewTokenBucketRateLimiterWithClock(qps float32, burst int, c Clock) flowcontrol.RateLimiter {
	limiter := rate.NewLimiter(rate.Limit(qps), burst)
	return newTokenBucketRateLimiter(limiter, c, qps)
}

func newTokenBucketRateLimiter(limiter *rate.Limiter, c Clock, qps float32) flowcontrol.RateLimiter {
	return &tokenBucketRateLimiter{
		limiter: limiter,
		clock:   c,
		qps:     qps,
	}
}

func (t *tokenBucketRateLimiter) TryAccept() bool {
	return t.limiter.AllowN(t.clock.Now(), 1)
}

// Accept will block until a token becomes available
func (t *tokenBucketRateLimiter) Accept() {
	now := t.clock.Now()
	t.clock.Sleep(t.limiter.ReserveN(now, 1).DelayFrom(now))
}

func (t *tokenBucketRateLimiter) Stop() {
}

func (t *tokenBucketRateLimiter) QPS() float32 {
	return t.qps
}

func (t *tokenBucketRateLimiter) Wait(ctx context.Context) error {
	return t.limiter.Wait(ctx)
}

type fakeAlwaysRateLimiter struct{}

// NewFakeAlwaysRateLimiter 用于测试
func NewFakeAlwaysRateLimiter() flowcontrol.RateLimiter {
	return &fakeAlwaysRateLimiter{}
}

func (t *fakeAlwaysRateLimiter) TryAccept() bool {
	return true
}

func (t *fakeAlwaysRateLimiter) Stop() {}

func (t *fakeAlwaysRateLimiter) Accept() {}

func (t *fakeAlwaysRateLimiter) QPS() float32 {
	return 1
}

func (t *fakeAlwaysRateLimiter) Wait(ctx context.Context) error {
	return nil
}

type fakeNeverRateLimiter struct {
	wg sync.WaitGroup
}

// NewFakeNeverRateLimiter 用于测试
func NewFakeNeverRateLimiter() flowcontrol.RateLimiter {
	rl := fakeNeverRateLimiter{}
	rl.wg.Add(1)
	return &rl
}

func (t *fakeNeverRateLimiter) TryAccept() bool {
	return false
}

func (t *fakeNeverRateLimiter) Stop() {
	t.wg.Done()
}

func (t *fakeNeverRateLimiter) Accept() {
	t.wg.Wait()
}

func (t *fakeNeverRateLimiter) QPS() float32 {
	return 1
}

func (t *fakeNeverRateLimiter) Wait(ctx context.Context) error {
	return errors.New("can not be accept")
}
