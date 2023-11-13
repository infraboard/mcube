package application

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	return r, nil
}
