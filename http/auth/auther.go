package auth

import "net/http"

// Auther 设置受保护路由使用的认证器
type Auther interface {
	Auth(http.Header) (authInfo interface{}, err error)
}

// The AutherFunc type is an adapter to allow the use of
// ordinary functions as Auther handlers. If f is a function
// with the appropriate signature, AutherFunc(f) is a
// Handler that calls f.
type AutherFunc func(http.Header) (authInfo interface{}, err error)

// Auth calls auth.
func (f AutherFunc) Auth(h http.Header) (authInfo interface{}, err error) {
	return f(h)
}
