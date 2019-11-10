package auth

import "net/http"

// Auther 设置受保护路由使用的认证器
type Auther interface {
	Auth(http.Header) (authInfo interface{}, err error)
}

