package app

// AllTemplate todo
const AllTemplate = `package all

import (
	// 加载服务模块
	_ "{{.PKG}}/apps/example/impl"
	_ "{{.PKG}}/apps/example/http"
)
`
