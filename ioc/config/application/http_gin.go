package application

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/ioc"
)

func NewGinRouterBuilder() *GinRouterBuilder {
	return &GinRouterBuilder{
		conf: &BuildConfig{},
	}
}

type GinRouterBuilder struct {
	conf *BuildConfig
}

func (b *GinRouterBuilder) Config(c *BuildConfig) {
	b.conf = c
}

func (b *GinRouterBuilder) Build() (http.Handler, error) {
	r := gin.Default()

	// 装载Ioc路由之前
	if b.conf.BeforeLoad != nil {
		b.conf.BeforeLoad(r)
	}

	// 装置子服务路由
	ioc.LoadGinApi(App().HTTPPrefix(), r)

	// 装载Ioc路由之后
	if b.conf.AfterLoad != nil {
		b.conf.AfterLoad(r)
	}
	return r, nil
}
