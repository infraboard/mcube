package router

import (
	"net/http"

	"github.com/infraboard/mcube/logger"
)

// Router 路由
type Router interface {
	// 添加中间件
	Use(m Middleware)
	// 添加受认证保护的路由
	AddProtected(method, path string, h http.HandlerFunc)
	// 添加公开路由, 所有人都可以访问
	AddPublict(method, path string, h http.HandlerFunc)
	// 开始认证时 使用的认证器
	SetAuther(Auther)

	// 实现标准库路由
	ServeHTTP(http.ResponseWriter, *http.Request)

	// 获取当前的路由条目信息
	GetEndpoints() *EntrySet

	// EnableAPIRoot 将服务路由表通过路径/暴露出去
	EnableAPIRoot()

	// 设置路由的Logger, 用于Debug
	SetLogger(logger.Logger)
	// SetLabel 设置路由标签, 作用于Entry上
	SetLabel(...*Label)

	// 子路由
	SubRouter(basePath string) SubRouter
}

// SubRouter 子路由或者分组路由
type SubRouter interface {
	// 添加中间件
	Use(m Middleware)
	// 添加受认证保护的路由
	AddProtected(method, path string, h http.HandlerFunc)
	// 添加公开路由, 所有人都可以访问
	AddPublict(method, path string, h http.HandlerFunc)
	// With独立作用于某一个Handler
	With(m ...Middleware) SubRouter
	// SetLabel 设置子路由标签, 作用于Entry上
	SetLabel(...*Label)
	// ResourceRouter 资源路由器, 主要用于设置路由标签和资源名称,
	// 方便配置灵活的权限策略
	ResourceRouter(resourceName string, labels ...*Label) SubRouter
}
