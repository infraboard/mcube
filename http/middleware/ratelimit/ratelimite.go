package ratelimit

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/infraboard/mcube/flowcontrol"
	"github.com/infraboard/mcube/flowcontrol/tokenbucket"
	"github.com/infraboard/mcube/ioc/config/logger"
	"github.com/rs/zerolog"
)

// NewGlobalModeLimiter todo
func NewGlobalModeLimiter(rate float64, capacity int64) *Limiter {
	return new(rate, capacity, GlobalMode)
}

// NewRemoteIPModeLimiter todo
func NewRemoteIPModeLimiter(rate float64, capacity int64) *Limiter {
	return new(rate, capacity, RemoteIPMode)
}

// NewHeaderKeyModeLimiter todo
func NewHeaderKeyModeLimiter(rate float64, capacity int64, headerKey string) *Limiter {
	l := new(rate, capacity, HeaderKeyMode)
	l.headerKey = headerKey
	return l
}

// NewCookieKeyModeLimiter todo
func NewCookieKeyModeLimiter(rate float64, capacity int64, cookieKey string) *Limiter {
	l := new(rate, capacity, CookieKeyMode)
	l.cookieKey = cookieKey
	return l
}

// New returns a new Logger instance
func new(rate float64, capacity int64, mode Mode) *Limiter {
	return &Limiter{
		rate:              rate,
		capacity:          capacity,
		limiters:          make(map[string]flowcontrol.RateLimiter),
		l:                 logger.Sub("Rate Limiter"),
		mode:              mode,
		remoteIPHeaderKey: []string{"X-Forwarded-For", "X-Real-IP"},
		maxSize:           1000,
	}
}

// Limiter is the structure for request rate limiter
type Limiter struct {
	rate              float64
	capacity          int64
	limiters          map[string]flowcontrol.RateLimiter
	count             int64
	maxSize           int64
	mu                sync.RWMutex
	l                 *zerolog.Logger
	mode              Mode
	remoteIPHeaderKey []string
	headerKey         string
	cookieKey         string
}

// SetMaxSize 设置最大值
func (l *Limiter) SetMaxSize(max int64) {
	l.maxSize = max
}

func (l *Limiter) removeExpired() int64 {
	var removedCount int64

	// 清理3分钟内无访问的限制器
	for k, v := range l.limiters {
		if time.Since(v.LastTakeTime()) > 3*time.Minute {
			delete(l.limiters, k)
			l.count--
			removedCount++
		}
	}

	return removedCount
}

// SetRemoteIPHeader 设置RemoteIP的头
func (l *Limiter) SetRemoteIPHeader(headers []string) {
	l.remoteIPHeaderKey = headers
}

// GetLimiter 获取对应的限制器
func (l *Limiter) GetLimiter(r *http.Request) flowcontrol.RateLimiter {
	id := l.getLimiterID(r)

	// 获取Limiter
	l.mu.RLock()
	limiter, exists := l.limiters[id]
	l.mu.RUnlock()

	if !exists {
		limiter = l.addLimiter(id)
	}

	return limiter
}

func (l *Limiter) addLimiter(id string) flowcontrol.RateLimiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.count >= l.maxSize {
		l.removeExpired()
	}

	limiter := tokenbucket.NewBucketWithRate(l.rate, l.capacity)
	l.limiters[id] = limiter
	l.count++

	return limiter
}

func (l *Limiter) getLimiterID(r *http.Request) string {
	switch l.mode {
	case GlobalMode:
		return "*"
	case RemoteIPMode:
		return l.remoteIP(r)
	case HeaderKeyMode:
		return l.headerKeyValue(r)
	case CookieKeyMode:
		return l.cookieKeyValue(r)
	}

	return ""
}

func (l *Limiter) remoteIP(r *http.Request) string {
	for _, header := range l.remoteIPHeaderKey {
		ip := r.Header.Get(header)
		if ip != "" {
			return ip
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		l.l.Error().Msgf("get net remote ip error, %s", err)
		return ""
	}

	return ip
}

func (l *Limiter) headerKeyValue(r *http.Request) string {
	return r.Header.Get(l.headerKey)
}

func (l *Limiter) cookieKeyValue(r *http.Request) string {
	ck, err := r.Cookie(l.cookieKey)
	if err != nil {
		l.l.Error().Msgf("get cookie key value error, %s", err)
		return ""
	}

	return ck.Value
}

// Handler 实现中间件
func (l *Limiter) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if !l.GetLimiter(r).TakeOneAvailable() {
			http.Error(rw, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(rw, r)
	})
}
