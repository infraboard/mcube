package httprouter_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/infraboard/mcube/http/auth/mock"
	"github.com/infraboard/mcube/http/router/httprouter"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}

func TestBase(t *testing.T) {
	r := httprouter.New()
	r.AddPublict("GET", "/", IndexHandler)

	t.Log(r.GetEndpoints())
	go http.ListenAndServe("0.0.0.0:8000", r)

	<-time.After(time.Second * 10)
}

func TestBaseWithAuther(t *testing.T) {
	r := httprouter.New()
	// 设置mock
	r.SetAuther(mock.NewMockAuther())

	r.AddProtected("GET", "/", IndexHandler)

	t.Log(r.GetEndpoints())
	go http.ListenAndServe("0.0.0.0:8000", r)

	<-time.After(time.Second * 10)
}
