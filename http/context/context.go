package context

import (
	"context"
	"net/http"

	httppb "github.com/infraboard/mcube/pb/http"
	"github.com/julienschmidt/httprouter"
)

// context包的WithValue函数, 官方推荐key应该是非导出的, 非go内置的的类型
// 究其原因是: 当前ctx没有还会去父ctx获取，和JS原型链调用类似，如果在不知情的情况下用了和父ctx
// 相同的key并且是内置类型就可能会有异常情况
type key int

const defaultKey = key(1)

// WithContext 携带请求上下文
func WithContext(req *http.Request, rctx *ReqContext) *http.Request {
	ctx := context.WithValue(req.Context(), defaultKey, rctx)
	return req.WithContext(ctx)
}

// GetContext 获取请求上下文中的数据
func GetContext(req *http.Request) *ReqContext {
	if v, ok := req.Context().Value(defaultKey).(*ReqContext); ok {
		return v
	}

	return new(ReqContext)
}

// NewReqContext 创建请假上下文实例
func NewReqContext() *ReqContext {
	return &ReqContext{}
}

// ReqContext context
type ReqContext struct {
	Entry    *httppb.Entry
	PS       httprouter.Params
	AuthInfo interface{}
}
