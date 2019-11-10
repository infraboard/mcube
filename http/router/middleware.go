package router

import (
	"net/http"

	"github.com/infraboard/mcube/http/context"
	"github.com/infraboard/mcube/http/response"
)

// Middleware 中间件的函数签名
type Middleware interface {
	Wrap(http.Handler) http.Handler
}

// NewAutherMiddleware 初始化一个认证中间件
func NewAutherMiddleware(auther Auther) *AutherMiddleware {
	return &AutherMiddleware{
		auther: auther,
	}
}

// AutherMiddleware 认证中间件
type AutherMiddleware struct {
	auther Auther
}

// Wrap 实现中间件
func (m *AutherMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		// 使用auther进行认证
		authInfo, err := m.auther.Auth(r.Header)
		if err != nil {
			response.Failed(wr, err)
			return
		}

		rc := context.GetContext(r)
		rc.AuthInfo = authInfo
		context.WithContext(r, rc)

		// next handler
		next.ServeHTTP(wr, r)
	})
}
