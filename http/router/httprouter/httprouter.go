package httprouter

import (
	"log"
	"net/http"

	"github.com/infraboard/mcube/http/auth"
	"github.com/infraboard/mcube/http/context"
	"github.com/infraboard/mcube/http/response"
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
	mergedHandler   http.Handler
}

// New 基于社区的httprouter进行封装
func New() router.Router {
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

	r.mergedHandler = r.r
	for i := len(r.middlewareChain) - 1; i >= 0; i-- {
		r.mergedHandler = r.middlewareChain[i].Handler(r.mergedHandler)
	}
}

func (r *httpRouter) AddProtected(method, path string, h http.HandlerFunc) {
	e := &entry{
		Entry: router.Entry{
			Name:      router.GetHandlerFuncName(h),
			Method:    method,
			Path:      path,
			Protected: true,
		},
		h: h,
	}
	r.add(e)
}

func (r *httpRouter) AddPublict(method, path string, h http.HandlerFunc) {
	e := &entry{
		Entry: router.Entry{
			Name:      router.GetHandlerFuncName(h),
			Method:    method,
			Path:      path,
			Protected: false,
		},
		h: h,
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

// ServeHTTP 扩展http ResponseWriter 接口, 使用扩展后兼容的接口替换掉原来的reponse对象
// 方便后期对response做处理
func (r *httpRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r.mergedHandler != nil {
		r.mergedHandler.ServeHTTP(response.NewResponse(w), req)
		return
	}

	r.r.ServeHTTP(response.NewResponse(w), req)
}

func (r *httpRouter) GetEndpoints() []router.Entry {
	entries := make([]router.Entry, 0, len(r.entries))
	for i := range r.entries {
		entries = append(entries, r.entries[i].Entry)
	}

	return entries
}

func (r *httpRouter) EnableAPIRoot() {
	r.AddPublict("GET", "/", r.apiRoot)
}

func (r *httpRouter) apiRoot(w http.ResponseWriter, req *http.Request) {
	response.Success(w, r.entries)
	return
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
	var handler http.Handler

	if e.Protected && r.authMiddleware == nil {
		r.debug("add projected endpoint, but no auther set")
	}

	if e.Protected && r.authMiddleware != nil {
		handler = r.authMiddleware.Handler(e.h)
	} else {
		handler = http.HandlerFunc(e.h)
	}

	r.addHandler(e.Method, e.Path, handler)
	r.addEntry(e)
}

func (r *httpRouter) addHandler(method, path string, h http.Handler) {
	r.r.Handle(method, path,
		func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
			rc := &context.ReqContext{
				PS: ps,
			}
			req = context.WithContext(req, rc)
			h.ServeHTTP(w, req)

		},
	)
}

func (r *httpRouter) addEntry(e *entry) {
	r.entries = append(r.entries, e)
}
