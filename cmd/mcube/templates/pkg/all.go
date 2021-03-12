package pkg

// AllTemplate todo
const AllTemplate = `package all

import (
	// 加载服务模块
	_ "{{.PKG}}/pkg/example/impl"
	_ "{{.PKG}}/pkg/example/http"
)
`
