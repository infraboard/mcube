package router

import (
	"net/http"

	"github.com/infraboard/mcube/http/auth"
	"github.com/infraboard/mcube/logger"
)

// Entry 路由条目
type Entry struct {
	Name   string
	Method string
	Path   string
	Desc   string
	Tags   map[string]string
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
	SetAuther(auth.Auther)

	// 实现标准库路由
	ServeHTTP(http.ResponseWriter, *http.Request)

	// 获取当前的路由条目信息
	GetEndpoints() []Entry

	// 设置路由的Logger, 用于Debug
	SetLogger(logger.Logger)

	// 生成子路由实例
	SubRouter(namespace string) SubRouter
}

// SubRouter 子路由或者分组路由
type SubRouter interface {
	// 添加中间件
	Use(m Middleware)
	// With独立作用于某一个Handler
	With(m Middleware) SubRouter
	// 添加受认证保护的路由
	AddProtected(method, path string, h http.HandlerFunc)
	// 添加公开路由, 所有人都可以访问
	AddPublict(method, path string, h http.HandlerFunc)
}
