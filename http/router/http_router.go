package router

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/infraboard/mcube/http/auth"
	"github.com/infraboard/mcube/http/context"
	"github.com/julienschmidt/httprouter"
)

// Entry 路由条目
type entry struct {
	Entry

	needAuth bool
	h        http.HandlerFunc
}

type httpRouter struct {
	middlewareChain []Middleware
	entries         []*entry
	authMiddleware  Middleware
	r               *httprouter.Router
}

// NewHTTPRouter 基于社区的httprouter进行封装
func NewHTTPRouter() Router {
	r := &httpRouter{
		r: &httprouter.Router{
			RedirectTrailingSlash:  true,
			RedirectFixedPath:      true,
			HandleMethodNotAllowed: true,
		},
	}

	return r
}

func (r *httpRouter) Use(m Middleware) {
	r.middlewareChain = append(r.middlewareChain, m)
}

func (r *httpRouter) AddProtected(method, path string, h http.HandlerFunc) {
	e := &entry{
		Entry: Entry{
			Name:   r.getHandlerName(h),
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
		Entry: Entry{
			Name:   r.getHandlerName(h),
			Method: method,
			Path:   path,
		},
		needAuth: false,
		h:        h,
	}
	r.add(e)
}

func (r *httpRouter) SetAuther(at auth.Auther) {
	am := newAutherMiddleware(at)
	r.authMiddleware = am
	return
}

// ServeHTTP 交给httprouter处理
func (r *httpRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.r.ServeHTTP(w, req)
}

func (r *httpRouter) GetEndpoints() []Entry {
	entries := make([]Entry, 0, len(r.entries))
	for i := range r.entries {
		entries = append(entries, r.entries[i].Entry)
	}

	return entries
}

func (r *httpRouter) add(e *entry) {
	var mergedHandler http.Handler

	if e.needAuth {
		r.Use(r.authMiddleware)
	}

	mergedHandler = http.HandlerFunc(e.h)
	for i := len(r.middlewareChain) - 1; i >= 0; i-- {
		mergedHandler = r.middlewareChain[i].Wrap(mergedHandler)
	}

	r.r.Handle(e.Method, e.Path,
		func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
			rc := &context.ReqContext{
				PS: ps,
			}
			context.WithContext(req, rc)
			mergedHandler.ServeHTTP(w, req)
		},
	)
	r.entries = append(r.entries, e)
}

func (r *httpRouter) getHandlerName(h http.Handler, seps ...rune) string {
	// 获取函数名称
	fn := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()

	// 默认使用.分隔
	if len(seps) == 0 {
		seps = append(seps, '.')
	}

	// 用 seps 进行分割
	fields := strings.FieldsFunc(fn, func(sep rune) bool {
		for _, s := range seps {
			if sep == s {
				return true
			}
		}
		return false
	})

	if size := len(fields); size > 0 {
		return fields[size-1]
	}

	return ""
}
