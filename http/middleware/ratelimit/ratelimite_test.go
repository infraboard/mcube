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

func TestRateLimte(t *testing.T) {
	should := require.New(t)

	router := httprouter.New()
	router.Use(ratelimit.NewALLModeLimiter(10, 10))
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

func init() {
	zap.DevelopmentSetup()
	zap.L()
}
