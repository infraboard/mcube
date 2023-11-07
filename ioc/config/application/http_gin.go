package application

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewGinRouterBuilder() *GinRouterBuilder {
	return &GinRouterBuilder{}
}

type GinRouterBuilder struct {
	Router *gin.Engine
}

func (b *GinRouterBuilder) Init() error {
	b.Router = gin.Default()
	return nil
}

func (b *GinRouterBuilder) GetRouter() http.Handler {
	return b.Router
}
