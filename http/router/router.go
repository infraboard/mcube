package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/http/context"
	"github.com/infraboard/mcube/http/response"
)

// MyRouter is an hack for httprouter
type MyRouter struct {
	Router *httprouter.Router

	urlPrefix   string
	v1endpoints map[string]map[string]string
}

// NewRouter use to new an router
func NewRouter() *MyRouter {
	hrouter := &httprouter.Router{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
	}

	ep := make(map[string]map[string]string)

	r := &MyRouter{Router: hrouter, v1endpoints: ep}
	return r
}

// SetURLPrefix 设置路由前缀
func (r *MyRouter) SetURLPrefix(prefix string) {
	r.urlPrefix = prefix
}

// handlerWithAuth is an adapter which allows the usage of an http.Handler as a
// request handle.
func (r *MyRouter) handlerWithAuth(method, path, featureName string, handler http.Handler) {
	r.Router.Handle(method, path,
		func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
			rctx := new(context.ReqContext)
			authHeader := req.Header.Get("Authorization")
			if authHeader == "" {
				response.Failed(w, exception.NewAPIException(exception.Unauthorized,
					"Authorization missed in header"))
				return
			}

			headerSlice := strings.Split(authHeader, " ")
			if len(headerSlice) != 2 {
				response.Failed(w, exception.NewAPIException(exception.TokenFormatInvalidate,
					"Authorization header value is not validated, must be: {token_type} {token}"))
				return
			}

			access := headerSlice[1]

			fmt.Printf(access)
			// t, err := dao.Auth.ValidateToken(access, featureName)
			// if err != nil {
			// 	response.Failed(w, err)
			// 	return
			// }

			// rctx.Token = t

			rctx.PS = ps
			context.SetReqContext(req, w, handler, rctx)
		},
	)
}

// handlerWithPS is an adapter which allows the usage of an http.Handler as a
// request handle.
func (r *MyRouter) handlerWithPS(method, path, featureName string, handler http.Handler) {
	r.Router.Handle(method, path,
		func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
			rctx := new(context.ReqContext)
			rctx.PS = ps
			context.SetReqContext(req, w, handler, rctx)
		},
	)
}

// ProtectedHandlerFunc 需要通过keyauth系统进行权限认证
func (r *MyRouter) ProtectedHandlerFunc(method, path, featureName string, handleFunc http.HandlerFunc) {
	path = r.urlPrefix + path

	_, ok := r.v1endpoints[method]
	if !ok {
		mm := make(map[string]string)
		r.v1endpoints[method] = mm
	}
	if featureName != "" {
		r.v1endpoints[method][featureName] = path
	}

	r.handlerWithAuth(method, path, featureName, http.HandlerFunc(handleFunc))
}

// PublicHandlerFunc 无须认证, 公开给所有人访问
func (r *MyRouter) PublicHandlerFunc(method, path, featureName string, handleFunc http.HandlerFunc) {
	path = r.urlPrefix + path

	switch featureName {
	case "":
	default:
		_, ok := r.v1endpoints[method]
		if !ok {
			mm := make(map[string]string)
			r.v1endpoints[method] = mm
		}

		if featureName != "" {
			r.v1endpoints[method][featureName] = path
		}
	}

	r.handlerWithPS(method, path, featureName, http.HandlerFunc(handleFunc))
}

// AddVRootIndex add root to api
func (r *MyRouter) AddRootIndex() {
	r.Router.Handle("GET", "/",
		func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
			response.Success(w, http.StatusOK, r.v1endpoints)
			return
		},
	)
}

// GetEndpoints get router's fl
func (r *MyRouter) GetEndpoints() map[string]map[string]string {
	return r.v1endpoints
}
