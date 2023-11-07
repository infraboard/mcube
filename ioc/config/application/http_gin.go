package application

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

func NewGinRouterBuilder() *GinRouterBuilder {
	return &GinRouterBuilder{}
}

type GinRouterBuilder struct {
	lock   sync.Mutex
	Router *gin.Engine
}

func (b *GinRouterBuilder) Init() error {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.Router == nil {
		b.Router = gin.Default()
	}
	return nil
}

func (b *GinRouterBuilder) GetRouter() http.Handler {
	return b.Router
}
