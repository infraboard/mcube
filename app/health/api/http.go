package api

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/app/health"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	h = &handler{}
)

type handler struct {
	service healthgrpc.HealthServer
	log     logger.Logger
}

func (h *handler) Config() error {
	h.log = zap.L().Named(h.Name())
	h.service = app.GetInternalApp(health.AppName).(healthgrpc.HealthServer)
	return nil
}

func (h *handler) Name() string {
	return health.AppName
}

func (h *handler) Version() string {
	return "v1"
}

func (h *handler) Registry(ws *restful.WebService) {
	tags := []string{"健康检查"}
	ws.Route(ws.GET("/").To(h.Check).
		Doc("查询服务当前状态").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", Health{}))
}

func init() {
	app.RegistryRESTfulApp(h)
}
