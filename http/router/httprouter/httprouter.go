package httprouter

import (
	"log"
	"net/http"
	"strings"

	"github.com/infraboard/mcube/http/context"
	"github.com/infraboard/mcube/http/response"
	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/logger"
	"github.com/julienschmidt/httprouter"
)

type httpRouter struct {
	r *httprouter.Router
	l logger.Logger

	middlewareChain []router.Middleware
	entrySet        *entrySet
	auther          router.Auther
	mergedHandler   http.Handler
	labels          []*router.Label
	notFound        http.Handler
}

// New 基于社区的httprouter进行封装
func New() router.Router {
	r := &httpRouter{
		notFound: http.HandlerFunc(http.NotFound),
		entrySet: newEntrySet(),
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
		Entry: &router.Entry{
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
		Entry: &router.Entry{
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

func (r *httpRouter) GetEndpoints() []router.Entry {
	return r.entrySet.ShowEntries()
}

func (r *httpRouter) EnableAPIRoot() {
	r.AddPublict("GET", "/", r.apiRoot)
}

func (r *httpRouter) apiRoot(w http.ResponseWriter, req *http.Request) {
	response.Success(w, r.entrySet.ShowEntries)
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

// 在添加路由时 装饰了认证逻辑
func (r *httpRouter) addHandler(protected bool, method, path string, h http.Handler) {
	wrapper := func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		// 使用auther进行认证
		if protected && r.auther != nil {
			entry := r.findEntry(method, path)
			if entry == nil {
				r.notFound.ServeHTTP(w, req)
				return
			}

			authInfo, err := r.auther.Auth(req.Header, *entry.Entry)
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
	}
	r.r.Handle(method, path, wrapper)
}

func (r *httpRouter) addEntry(e *entry) {
	// 主路由全局标签
	for i := range r.labels {
		kv := r.labels[i]
		e.Labels[kv.Key()] = kv.Value()
	}

	if err := r.entrySet.AddEntry(e); err != nil {
		panic(err)
	}
}

func (r *httpRouter) findEntry(method, path string) *entry {
	hander, paras, _ := r.r.Lookup(method, path)
	if hander == nil {
		return nil
	}

	targetPath := path
	for _, kv := range paras {
		targetPath = strings.Replace(targetPath, ":"+kv.Key, kv.Value, 1)
	}

	return r.entrySet.FindEntry(targetPath, method)

}
