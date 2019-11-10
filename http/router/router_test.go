package router_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/infraboard/mcube/http/auth/mock"
	"github.com/infraboard/mcube/http/router"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}

func TestBase(t *testing.T) {
	r := router.NewHTTPRouter()
	r.AddPublict("GET", "/", IndexHandler)

	t.Log(r.GetEndpoints())
	go http.ListenAndServe("0.0.0.0:8000", r)

	<-time.After(time.Second * 10)
}

func TestBaseWithAuther(t *testing.T) {
	r := router.NewHTTPRouter()
	// 设置mock
	r.SetAuther(mock.NewMockAuther())

	r.AddProtected("GET", "/", IndexHandler)

	t.Log(r.GetEndpoints())
	go http.ListenAndServe("0.0.0.0:8000", r)

	<-time.After(time.Second * 10)
}
