package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/infraboard/mcube/http/middleware/accesslog"
	"github.com/infraboard/mcube/http/router/httprouter"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "hello world")
}

func main() {
	router := httprouter.NewHTTPRouter()
	lm := accesslog.NewLogger()
	router.Use(lm)
	router.AddPublict("GET", "/", IndexHandler)

	go http.ListenAndServe("0.0.0.0:8000", router)

	<-time.After(time.Second * 10)
}
