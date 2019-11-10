package router

import (
	"net/http"
)

// Entry 路由条目
type Entry struct {
	Name   string
	Method string
	Path   string
	Desc   string
	Tags   map[string]string
}

// Auther 设置受保护路由使用的认证器
type Auther interface {
	Auth(http.Header) (authInfo interface{}, err error)
}

// Router 路由
type Router interface {
	// 添加中间件
	Use(m Middleware)
	// 添加受认证保护的路由
	AddProtected(method, path string, h http.HandlerFunc)
	// 添加公开路由, 所有人都可以访问
	AddPublict(method, path string, h http.HandlerFunc)
	// 开始认证时 使用的认证器
	SetAuther(Auther) error

	ServeHTTP(http.ResponseWriter, *http.Request)

	GetEndpoints() []Entry
}
