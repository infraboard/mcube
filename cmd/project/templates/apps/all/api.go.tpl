package all

import (
	// 注册所有HTTP服务模块, 暴露给框架HTTP服务器加载
	{{ if .GenExample -}}
	_ "{{.PKG}}/apps/book/api"
	{{- end }}
)