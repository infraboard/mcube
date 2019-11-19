package main

import (
	"net/http"
	"time"

	"github.com/infraboard/mcube/http/middleware/accesslog"
	"github.com/infraboard/mcube/http/middleware/cors"
	"github.com/infraboard/mcube/http/middleware/recovery"
	"github.com/infraboard/mcube/http/router/httprouter"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	panic("test")
}

func main() {
	router := httprouter.NewHTTPRouter()

	lm := accesslog.New()
	router.Use(lm)

	rm := recovery.New()
	router.Use(rm)

	cm := cors.AllowAll()
	router.Use(cm)

	router.AddPublict("GET", "/", indexHandler)

	go http.ListenAndServe("0.0.0.0:8000", router)

	<-time.After(time.Second * 10)
}
