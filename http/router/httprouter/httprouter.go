package httprouter

import (
	"log"
	"net/http"

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
	r *httprouter.Router
	l logger.Logger

	middlewareChain []router.Middleware
	entries         []*entry
	auther          router.Auther
	mergedHandler   http.Handler
	labels          []*router.Label
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
			Labels:    map[string]string{},
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

func (r *httpRouter) SetAuther(at router.Auther) {
	r.auther = at
	return
}

func (r *httpRouter) SetLogger(logger logger.Logger) {
	r.l = logger
}

func (r *httpRouter) SetLabel(labels ...*router.Label) {
	r.labels = append(r.labels, labels...)
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

func (r *httpRouter) GetEndpoints() *router.EntrySet {
	es := router.NewEntrySet()

	for i := range r.entries {
		es.AddEntry(r.entries[i].Entry)
	}

	return es
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

	if e.Protected && r.auther == nil {
		r.debug("add projected endpoint, but no auther set")
	}

	handler = http.HandlerFunc(e.h)

	r.addHandler(e.Protected, e.Method, e.Path, handler)
	r.addEntry(e)
}

func (r *httpRouter) addHandler(protected bool, method, path string, h http.Handler) {
	r.r.Handle(method, path,
		func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
			// 使用auther进行认证
			if protected && r.auther != nil {
				authInfo, err := r.auther.Auth(req.Header)
				if err != nil {
					response.Failed(w, err)
					return
				}

				rc := context.GetContext(req)
				rc.AuthInfo = authInfo
				req = context.WithContext(req, rc)
			}

			// 未包含和auth为nil
			if !protected || r.auther == nil {
				rc := &context.ReqContext{
					PS: ps,
				}

				req = context.WithContext(req, rc)
			}

			h.ServeHTTP(w, req)
		},
	)
}

func (r *httpRouter) addEntry(e *entry) {
	// 主路由全局标签
	for i := range r.labels {
		kv := r.labels[i]
		e.Labels[kv.Key()] = kv.Value()
	}
	r.entries = append(r.entries, e)
}
