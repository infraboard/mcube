package httprouter

import (
	"fmt"
	"net/http"
	"path"

	"github.com/infraboard/mcube/http/router"
	httppb "github.com/infraboard/mcube/pb/http"
)

func newSubRouter(basePath string, root *httpRouter) *subRouter {
	return &subRouter{
		basePath:          basePath,
		root:              root,
		authEnable:        root.authEnable,
		permissionEnable:  root.permissionEnable,
		auditLog:          root.auditLog,
		requiredNamespace: root.requiredNamespace,
	}
}

type subRouter struct {
	basePath          string
	resourceName      string
	resourceBasePath  string
	root              *httpRouter
	labels            []*httppb.Label
	middlewareChain   []router.Middleware
	authEnable        bool
	permissionEnable  bool
	allow             []string
	auditLog          bool
	requiredNamespace bool
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

func (r *subRouter) Handle(method, path string, h http.HandlerFunc) httppb.EntryDecorator {
	e := &entry{
		Entry: &httppb.Entry{
			Resource:          r.resourceName,
			Method:            method,
			Path:              path,
			FunctionName:      router.GetHandlerFuncName(h),
			Labels:            map[string]string{},
			AuthEnable:        r.authEnable,
			PermissionEnable:  r.permissionEnable,
			AuditLog:          r.auditLog,
			Allow:             r.allow,
			RequiredNamespace: r.requiredNamespace,
		},
		h: h,
	}

	r.add(e)
	return e.Entry
}

func (r *subRouter) SetLabel(labels ...*httppb.Label) {
	r.labels = append(r.labels, labels...)
}

func (r *subRouter) Auth(isEnable bool) {
	r.authEnable = isEnable
}

func (r *subRouter) Permission(isEnable bool) {
	r.permissionEnable = isEnable
}

func (r *subRouter) Allow(targets ...fmt.Stringer) {
	for i := range targets {
		r.allow = append(r.allow, targets[i].String())
	}
}

func (r *subRouter) AuditLog(isEnable bool) {
	r.auditLog = isEnable
}

func (r *subRouter) RequiredNamespace(isEnable bool) {
	r.requiredNamespace = isEnable
}

func (r *subRouter) ResourceRouter(resourceName string, labels ...*httppb.Label) router.ResourceRouter {
	return &subRouter{
		resourceName:      resourceName,
		resourceBasePath:  r.basePath,
		basePath:          r.basePath,
		root:              r.root,
		labels:            append(r.labels, labels...),
		authEnable:        r.authEnable,
		permissionEnable:  r.permissionEnable,
		auditLog:          r.auditLog,
		allow:             r.allow,
		requiredNamespace: r.requiredNamespace,
	}
}

func (r *subRouter) BasePath(path string) {
	r.basePath = r.resourceBasePath + "/" + path
}

func (r *subRouter) add(e *entry) {
	// 子路由全局标签
	for i := range r.labels {
		kv := r.labels[i]
		e.Labels[kv.Key()] = kv.Value()
	}

	mergedHandler := r.combineHandler(e)
	e.Path = r.calculateAbsolutePath(e.Path)
	r.root.addHandler(e.Method, e.Path, mergedHandler)
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
