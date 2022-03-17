package router

import (
	"fmt"
	"net/http"

	"github.com/infraboard/mcube/logger"
	httppb "github.com/infraboard/mcube/pb/http"
)

// Router 路由
type Router interface {
	// 添加中间件
	Use(m Middleware)

	// 是否启用用户身份验证
	Auth(isEnable bool)

	// 是否启用用户权限验证
	Permission(isEnable bool)

	// 允许target标识
	Allow(targets ...fmt.Stringer)

	// 是否开启审计日志
	AuditLog(isEnable bool)

	// 是否需要NameSpace
	RequiredNamespace(isEnable bool)

	// 添加受认证保护的路由
	Handle(method, path string, h http.HandlerFunc) httppb.EntryDecorator

	// 开始认证时 使用的认证器
	SetAuther(Auther)

	// 开始审计器, 记录用户操作
	SetAuditer(Auditer)

	// 实现标准库路由
	ServeHTTP(http.ResponseWriter, *http.Request)

	// 获取当前的路由条目信息
	GetEndpoints() *httppb.EntrySet

	// EnableAPIRoot 将服务路由表通过路径/暴露出去
	EnableAPIRoot()

	// 设置路由的Logger, 用于Debug
	SetLogger(logger.Logger)

	// SetLabel 设置路由标签, 作用于Entry上
	SetLabel(...*httppb.Label)

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

	// 允许target标识
	Allow(targets ...fmt.Stringer)

	// 是否开启审计日志
	AuditLog(isEnable bool)

	// 是否需要NameSpace
	RequiredNamespace(isEnable bool)

	// 添加中间件
	Use(m Middleware)

	// SetLabel 设置路由标签, 作用于Entry上
	SetLabel(...*httppb.Label)

	// With独立作用于某一个Handler
	With(m ...Middleware) SubRouter

	// 添加受认证保护的路由
	Handle(method, path string, h http.HandlerFunc) httppb.EntryDecorator

	// ResourceRouter 资源路由器, 主要用于设置路由标签和资源名称,方便配置灵活的权限策略
	ResourceRouter(resourceName string, labels ...*httppb.Label) ResourceRouter
}
