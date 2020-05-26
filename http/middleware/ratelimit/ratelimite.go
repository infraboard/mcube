package ratelimit

import (
	"net"
	"net/http"

	"github.com/infraboard/mcube/flowcontrol"
	"github.com/infraboard/mcube/flowcontrol/tokenbucket"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
)

// NewALLModeLimiter todo
func NewALLModeLimiter(rate float64, capacity int64) *Limiter {
	return new(rate, capacity, ALLMode)
}

// New returns a new Logger instance
func new(rate float64, capacity int64, mode Mode) *Limiter {
	return &Limiter{
		limiter:           tokenbucket.NewBucketWithRate(rate, capacity),
		l:                 zap.L().Named("Rate Limiter"),
		mode:              mode,
		remoteIPHeaderKey: []string{"X-Forwarded-For", "X-Real-IP"},
	}
}

// Limiter is the structure for request rate limiter
type Limiter struct {
	limiter           flowcontrol.RateLimiter
	l                 logger.Logger
	mode              Mode
	remoteIPHeaderKey []string
	headerKey         string
	cookieKey         string
}

// SetRemoteIPHeader 设置RemoteIP的头
func (l *Limiter) SetRemoteIPHeader(headers []string) {
	l.remoteIPHeaderKey = headers
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
		l.l.Errorf("get net remote ip error, %s", err)
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
		l.l.Errorf("get cookie key value error, %s", err)
		return ""
	}

	return ck.Value
}

// Handler 实现中间件
func (l *Limiter) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if !l.limiter.TakeAvailableOnce() {
			http.Error(rw, http.StatusText(429), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(rw, r)
	})
}
