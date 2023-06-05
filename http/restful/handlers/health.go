package handlers

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/app/health"
	"github.com/infraboard/mcube/http/restful/response"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		service:         app.GetInternalApp(health.AppName).(healthgrpc.HealthServer),
		log:             zap.L().Named("health_check"),
		HealthCheckPath: "/healthz",
	}
}

type HealthChecker struct {
	service         healthgrpc.HealthServer
	log             logger.Logger
	HealthCheckPath string
}

func (h *HealthChecker) WebService() *restful.WebService {
	ws := new(restful.WebService)
	tags := []string{"健康检查"}
	ws.Route(ws.GET(h.HealthCheckPath).To(h.Check).
		Doc("查询服务当前状态").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", Health{}))

	return ws
}

func (h *HealthChecker) Check(r *restful.Request, w *restful.Response) {
	req := health.NewHealthCheckRequest()
	resp, err := h.service.Check(
		r.Request.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, NewHealth(resp))
}

func NewHealth(hc *healthgrpc.HealthCheckResponse) *Health {
	return &Health{
		Status: hc.Status.String(),
	}
}

type Health struct {
	Status string `json:"status"`
}
