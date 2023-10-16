package api

import (
	"github.com/infraboard/mcube/http/label"
	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"{{.PKG}}/apps/book"
	"github.com/infraboard/mcube/app"
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

func (h *handler) Registry(r router.SubRouter) {
	rr := r.ResourceRouter("books")

	rr.BasePath("books")
	rr.Handle("POST", "/", h.CreateBook).AddLabel(label.Create)
	rr.Handle("GET", "/", h.QueryBook).AddLabel(label.List)
	rr.Handle("GET", "/:id", h.DescribeBook).AddLabel(label.Get)
	rr.Handle("PUT", "/:id", h.PutBook).AddLabel(label.Update)
	rr.Handle("PATCH", "/:id", h.PatchBook).AddLabel(label.Update)
	rr.Handle("DELETE", "/:id", h.DeleteBook).AddLabel(label.Delete)
}

func init() {
	app.RegistryHttpApp(h)
}
