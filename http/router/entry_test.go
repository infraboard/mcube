package router_test

import (
	"net/http"

	"github.com/infraboard/mcube/http/response"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	response.Success(w, "ok")
	return
}
