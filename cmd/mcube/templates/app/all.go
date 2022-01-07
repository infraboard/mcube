package app

const GRPC_SERVICE_REGISTRY_Template = `package all

import (
	// 注册所有GRPC服务模块, 暴露给框架GRPC服务器加载, 注意 导入有先后顺序
	// _ "{{.PKG}}/apps/example/impl"
)
`

const HTTP_SERVICE_REGISTRY_Template = `package all

import (
	// 注册所有HTTP服务模块, 暴露给框架HTTP服务器加载
	// _ "{{.PKG}}/apps/example/http"
)
`

const INTERNAL_SERVICE_REGISTRY_Template = `package all

import (
	// 注册所有内部服务模块, 无须对外暴露的服务, 用于内部依赖
)
`
