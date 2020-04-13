package router

import (
	"net/http"

	"github.com/infraboard/mcube/logger"
)

// Router 路由
type Router interface {
	// 添加中间件
	Use(m Middleware)

	// 是否启用用户身份验证
	Auth(isEnable bool)

	// 是否启用用户权限验证
	Permission(isEnable bool)

	// 添加受认证保护的路由
	Handle(method, path string, h http.HandlerFunc) EntryDecorator

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

// ResourceRouter 资源路由
type ResourceRouter interface {
	SubRouter
	// BasePath 设置资源路由的基础路径
	BasePath(path string)
}

// SubRouter 子路由或者分组路由
type SubRouter interface {
	// 是否启用用户身份验证
	Auth(isEnable bool)
	// 是否启用用户权限验证
	Permission(isEnable bool)
	// 添加中间件
	Use(m Middleware)
	// SetLabel 设置路由标签, 作用于Entry上
	SetLabel(...*Label)
	// With独立作用于某一个Handler
	With(m ...Middleware) SubRouter
	// 添加受认证保护的路由
	Handle(method, path string, h http.HandlerFunc) EntryDecorator
	// ResourceRouter 资源路由器, 主要用于设置路由标签和资源名称,方便配置灵活的权限策略
	ResourceRouter(resourceName string, labels ...*Label) ResourceRouter
}

// EntryDecorator 装饰
type EntryDecorator interface {
	// SetLabel 设置子路由标签, 作用于Entry上
	AddLabel(...*Label) EntryDecorator
	EnableAuth() EntryDecorator
	DisableAuth() EntryDecorator
	EnablePermission() EntryDecorator
	DisablePermission() EntryDecorator
}
