package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/v2/ioc"
)

func NewGinRouterBuilder() *GinRouterBuilder {
	return &GinRouterBuilder{
		RouterBuilderHooks: NewRouterBuilderHooks(),
	}
}

type GinRouterBuilder struct {
	*RouterBuilderHooks
}

func (b *GinRouterBuilder) Build() (http.Handler, error) {
	r := gin.Default()

	// 装载Ioc路由之前
	for _, fn := range b.Before {
		fn(r)
	}

	// 装置子服务路由
	ioc.LoadGinApi(Get().HTTPPrefix(), r)

	// 装载Ioc路由之后
	for _, fn := range b.After {
		fn(r)
	}

	return r, nil
}
