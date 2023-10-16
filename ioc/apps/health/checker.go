package health

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/http/restful/response"
	"github.com/infraboard/mcube/ioc/config/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

func NewDefaultHealthChecker() *HealthChecker {
	return NewHealthChecker(health.NewServer())
}

func NewHealthChecker(checker healthgrpc.HealthServer) *HealthChecker {
	return &HealthChecker{
		service:         checker,
		log:             logger.Sub("health_check"),
		HealthCheckPath: "/healthz",
	}
}

type HealthChecker struct {
	service         healthgrpc.HealthServer
	log             *zerolog.Logger
	HealthCheckPath string
}

func (h *HealthChecker) WebService() *restful.WebService {
	ws := new(restful.WebService)
	tags := []string{"健康检查"}
	ws.Route(ws.GET(h.HealthCheckPath).To(h.Check).
		Doc("查询服务当前状态").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", HealthCheckResponse{}))

	return ws
}

func (h *HealthChecker) Check(r *restful.Request, w *restful.Response) {
	req := NewHealthCheckRequest()
	resp, err := h.service.Check(
		r.Request.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}

	err = w.WriteAsJson(NewHealth(resp))
	if err != nil {
		zap.L().Errorf("send success response error, %s", err)
	}
}

func NewHealthCheckRequest() *healthgrpc.HealthCheckRequest {
	return &healthgrpc.HealthCheckRequest{}
}

func NewHealth(hc *healthgrpc.HealthCheckResponse) *HealthCheckResponse {
	return &HealthCheckResponse{
		Status: hc.Status.String(),
	}
}

type HealthCheckResponse struct {
	Status string `json:"status"`
}
