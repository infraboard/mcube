package ratelimit_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/infraboard/mcube/http/middleware/ratelimit"
	"github.com/infraboard/mcube/http/router/httprouter"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/stretchr/testify/require"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "hello world!")
}

func TestGlobalLimter(t *testing.T) {
	should := require.New(t)

	router := httprouter.New()
	router.Use(ratelimit.NewGlobalModeLimiter(10, 10))
	router.Handle("GET", "/", indexHandler)

	req, _ := http.NewRequest("GET", "/", nil)

	for i := 0; i < 12; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if i == 10 {
			should.Equal(429, w.Code)
		}
	}
}

func TestRemoteIPLimiter(t *testing.T) {
	should := require.New(t)

	router := httprouter.New()
	router.Use(ratelimit.NewRemoteIPModeLimiter(10, 10))
	router.Handle("GET", "/", indexHandler)

	req, _ := http.NewRequest("GET", "/", nil)

	for i := 0; i < 12; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if i == 10 {
			should.Equal(429, w.Code)
		}
	}

	for i := 0; i < 12; i++ {
		req.Header.Set("X-Real-IP", fmt.Sprintf("10.10.10.%d", i))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if i == 10 {
			should.Equal(200, w.Code)
		}
	}
}

func TestHeaderKeyLimiter(t *testing.T) {
	should := require.New(t)

	router := httprouter.New()
	router.Use(ratelimit.NewHeaderKeyModeLimiter(10, 10, "X-OAUTH-TOKEN"))
	router.Handle("GET", "/", indexHandler)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-OAUTH-TOKEN", "xxxx")

	for i := 0; i < 12; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if i == 10 {
			should.Equal(429, w.Code)
		}
	}

	for i := 0; i < 12; i++ {
		req.Header.Set("X-OAUTH-TOKEN", fmt.Sprintf("xxx.%d", i))

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if i == 10 {
			should.Equal(200, w.Code)
		}
	}
}

func TestCookieKeyLimiter(t *testing.T) {
	should := require.New(t)

	router := httprouter.New()
	router.Use(ratelimit.NewCookieKeyModeLimiter(10, 10, "token"))
	router.Handle("GET", "/", indexHandler)

	req, _ := http.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "xxx"})

	for i := 0; i < 12; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if i == 10 {
			should.Equal(429, w.Code)
		}
	}

	for i := 0; i < 12; i++ {
		req, _ = http.NewRequest("GET", "/", nil)
		req.AddCookie(&http.Cookie{Name: "token", Value: fmt.Sprintf("xxx.%d", i)})
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if i == 10 {
			should.Equal(200, w.Code)
		}
	}
}

func TestLimiterClean(t *testing.T) {
	should := require.New(t)

	router := httprouter.New()
	m := ratelimit.NewRemoteIPModeLimiter(10, 10)
	m.SetMaxSize(10)
	router.Use(m)
	router.Handle("GET", "/", indexHandler)

	req, _ := http.NewRequest("GET", "/", nil)
	for i := 0; i < 12; i++ {
		req.Header.Set("X-Real-IP", fmt.Sprintf("10.10.10.%d", i))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if i == 10 {
			should.Equal(200, w.Code)
		}
	}
}

func init() {
	zap.DevelopmentSetup()
	zap.L()
}
