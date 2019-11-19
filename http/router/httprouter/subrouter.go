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
	root            *httpRouter
	middlewareChain []router.Middleware
}

func (r *subRouter) Use(m router.Middleware) {
	r.middlewareChain = append(r.middlewareChain, m)
}

func (r *subRouter) With(m router.Middleware) router.SubRouter {
	return &subRouter{
		basePath:        r.basePath,
		root:            r.root,
		middlewareChain: r.middlewareChain,
	}
}

func (r *subRouter) AddProtected(method, path string, h http.HandlerFunc) {
	e := &entry{
		Entry: router.Entry{
			Name:   router.GetHandlerFuncName(h),
			Method: method,
			Path:   path,
		},
		needAuth: true,
		h:        h,
	}

	r.add(e)
}

func (r *subRouter) AddPublict(method, path string, h http.HandlerFunc) {
	e := &entry{
		Entry: router.Entry{
			Name:   router.GetHandlerFuncName(h),
			Method: method,
			Path:   path,
		},
		needAuth: false,
		h:        h,
	}

	r.add(e)
}

func (r *subRouter) add(e *entry) {
	mergedHandler := r.combineHandler(e)
	r.root.addHandler(e.Method, r.calculateAbsolutePath(e.Path), mergedHandler)
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
