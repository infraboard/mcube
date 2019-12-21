package httprouter

import (
	"net/http"
	"path"

	"github.com/infraboard/mcube/http/router"
)

func newSubRouter(basePath string, root *httpRouter) *subRouter {
	return &subRouter{
		basePath: basePath,
		root:     root,
	}
}

type subRouter struct {
	basePath        string
	resourceName    string
	root            *httpRouter
	labels          []*router.Label
	middlewareChain []router.Middleware
}

func (r *subRouter) Use(m router.Middleware) {
	r.middlewareChain = append(r.middlewareChain, m)
}

func (r *subRouter) With(m ...router.Middleware) router.SubRouter {
	// 继承原有的中间件
	sr := &subRouter{
		basePath:        r.basePath,
		root:            r.root,
		middlewareChain: r.middlewareChain,
	}

	// 添加新中间件
	sr.middlewareChain = append(sr.middlewareChain, m...)
	return sr
}

func (r *subRouter) AddProtected(method, path string, h http.HandlerFunc) {
	e := &entry{
		Entry: &router.Entry{
			Resource: r.resourceName,
			Method:   method,
			Path:     path,
			Labels:   map[string]string{},
		},
		needAuth: true,
		h:        h,
	}

	r.add(e)
}

func (r *subRouter) AddPublict(method, path string, h http.HandlerFunc) {
	e := &entry{
		Entry: &router.Entry{
			Resource: r.resourceName,
			Method:   method,
			Path:     path,
			Labels:   map[string]string{},
		},
		needAuth: false,
		h:        h,
	}

	r.add(e)
}

func (r *subRouter) SetLabel(labels ...*router.Label) {
	r.labels = append(r.labels, labels...)
}

func (r *subRouter) ResourceRouter(resourceName string) router.SubRouter {
	return &subRouter{
		resourceName: resourceName,
		basePath:     r.basePath + "/" + resourceName,
		root:         r.root,
	}
}

func (r *subRouter) add(e *entry) {
	// 子路由全局标签
	for i := range r.labels {
		kv := r.labels[i]
		e.Labels[kv.Key()] = kv.Value()
	}

	mergedHandler := r.combineHandler(e)
	e.Path = r.calculateAbsolutePath(e.Path)
	r.root.addHandler(e.Protected, e.Method, e.Path, mergedHandler)
	r.root.addEntry(e)
}

func (r *subRouter) combineHandler(e *entry) http.Handler {
	var mergedHandler http.Handler

	mergedHandler = http.HandlerFunc(e.h)
	for i := len(r.middlewareChain) - 1; i >= 0; i-- {
		mergedHandler = r.middlewareChain[i].Handler(mergedHandler)
	}

	return mergedHandler
}

func (r *subRouter) calculateAbsolutePath(relativePath string) string {
	if relativePath == "" {
		return r.basePath
	}

	finalPath := path.Join(r.basePath, relativePath)
	appendSlash := lastChar(relativePath) == '/' && lastChar(finalPath) != '/'
	if appendSlash {
		return finalPath + "/"
	}

	return finalPath
}

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}
