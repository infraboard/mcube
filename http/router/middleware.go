package router

import (
	"net/http"
)

// Middleware 中间件的函数签名
type Middleware interface {
	Handler(http.Handler) http.Handler
}

// MiddlewareFunc is an adapter to allow the use of ordinary functions as Negroni handlers.
// If f is a function with the appropriate signature, HandlerFunc(f) is a Handler object that calls f.
type MiddlewareFunc func(http.Handler) http.Handler

// Handler wrappe for function
func (h MiddlewareFunc) Handler(next http.Handler) http.Handler {
	return h(next)
}
