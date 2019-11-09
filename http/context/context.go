package context

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// context包的WithValue函数, 官方推荐key应该是非导出的, 非go内置的的类型
// 究其原因是: 当前ctx没有还会去父ctx获取，和JS原型链调用类似，如果在不知情的情况下用了和父ctx
// 相同的key并且是内置类型就可能会有异常情况
type key int

const defaultKey = key(1)

// ReqContext context
type ReqContext struct {
	PS httprouter.Params
}

// SetReqContext use to get context
func SetReqContext(req *http.Request, w http.ResponseWriter, next http.Handler, rctx *ReqContext) {
	ctx := context.WithValue(req.Context(), defaultKey, rctx)
	req = req.WithContext(ctx)
	next.ServeHTTP(w, req)
}

// GetParamsFromContext use to get httprouter ps from context
func GetParamsFromContext(req *http.Request) httprouter.Params {
	return req.Context().Value(defaultKey).(*ReqContext).PS
}
