package router_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/test"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}

func TestBase(t *testing.T) {
	r := router.NewHTTPRouter()
	r.AddPublict("GET", "/", test.IndexHandler)

	fmt.Println(r.GetEndpoints())
	go http.ListenAndServe("0.0.0.0:8000", r)

	<-time.After(time.Second * 2)
}
