package api

import (
	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/infraboard/mcube/app"

	"{{.PKG}}/apps/book"
)

var (
	h = &handler{}
)

type handler struct {
	service book.ServiceServer
	log     logger.Logger
}

func (h *handler) Config() error {
	h.log = logger.Sub(book.AppName)
	h.service = app.GetGrpcApp(book.AppName).(book.ServiceServer)
	return nil
}

func (h *handler) Name() string {
	return book.AppName
}

func (h *handler) Version() string {
	return "v1"
}

func (h *handler) Registry(r gin.IRouter) {
	r.POST("/", h.CreateBook)
	r.GET("/", h.QueryBook)
	r.GET("/:id", h.DescribeBook)
	r.PUT("/:id", h.PutBook)
	r.PATCH("/:id", h.PatchBook)
	r.DELETE("/:id", h.DeleteBook)
}

func init() {
	app.RegistryGinApp(h)
}
