package httprouter

import (
	"log"
	"net/http"

	"github.com/infraboard/mcube/http/auth"
	"github.com/infraboard/mcube/http/context"
	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/logger"
	"github.com/julienschmidt/httprouter"
)

// Entry 路由条目
type entry struct {
	router.Entry

	needAuth bool
	h        http.HandlerFunc
}

type httpRouter struct {
	middlewareChain []router.Middleware
	entries         []*entry
	authMiddleware  router.Middleware
	r               *httprouter.Router
	l               logger.Logger
}

// NewHTTPRouter 基于社区的httprouter进行封装
func NewHTTPRouter() router.Router {
	r := &httpRouter{
		r: &httprouter.Router{
			RedirectTrailingSlash:  true,
			RedirectFixedPath:      true,
			HandleMethodNotAllowed: true,
		},
	}

	return r
}

func (r *httpRouter) Use(m router.Middleware) {
	r.middlewareChain = append(r.middlewareChain, m)
}

func (r *httpRouter) AddProtected(method, path string, h http.HandlerFunc) {
	e := &entry{
		Entry: router.Entry{
			Name:   router.GetHandlerName(h),
			Method: method,
			Path:   path,
		},
		needAuth: true,
		h:        h,
	}
	r.add(e)
}

func (r *httpRouter) AddPublict(method, path string, h http.HandlerFunc) {
	e := &entry{
		Entry: router.Entry{
			Name:   router.GetHandlerName(h),
			Method: method,
			Path:   path,
		},
		needAuth: false,
		h:        h,
	}
	r.add(e)
}

func (r *httpRouter) SetAuther(at auth.Auther) {
	am := router.NewAutherMiddleware(at)
	r.authMiddleware = am
	return
}

func (r *httpRouter) SetLogger(logger logger.Logger) {
	r.l = logger
}

// ServeHTTP 交给httprouter处理
func (r *httpRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.r.ServeHTTP(w, req)
}

func (r *httpRouter) GetEndpoints() []router.Entry {
	entries := make([]router.Entry, 0, len(r.entries))
	for i := range r.entries {
		entries = append(entries, r.entries[i].Entry)
	}

	return entries
}

func (r *httpRouter) SubRouter(basePath string) router.SubRouter {
	return newSubRouter(basePath, r)
}

func (r *httpRouter) debug(format string, args ...interface{}) {
	if r.l != nil {
		r.l.Debugf(format, args...)
		return
	}

	log.Printf(format, args...)
}

func (r *httpRouter) add(e *entry) {
	mergedHandler := r.combineHandler(e)
	r.addHandler(e.Method, e.Path, mergedHandler)
	r.addEntry(e)
}

func (r *httpRouter) combineHandler(e *entry) http.Handler {
	var mergedHandler http.Handler

	if e.needAuth && r.authMiddleware == nil {
		r.debug("add protect handFunc, please SetAuther first.")
	}

	if e.needAuth && r.authMiddleware != nil {
		r.Use(r.authMiddleware)
	}

	mergedHandler = http.HandlerFunc(e.h)
	for i := len(r.middlewareChain) - 1; i >= 0; i-- {
		mergedHandler = r.middlewareChain[i].Wrap(mergedHandler)
	}

	return mergedHandler
}

func (r *httpRouter) addHandler(method, path string, mergedHandler http.Handler) {
	r.r.Handle(method, path,
		func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
			rc := &context.ReqContext{
				PS: ps,
			}
			context.WithContext(req, rc)
			mergedHandler.ServeHTTP(w, req)
		},
	)
}

func (r *httpRouter) addEntry(e *entry) {
	r.entries = append(r.entries, e)
}
