package restful

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/http/restful/response"
	"github.com/infraboard/mcube/v2/ioc"
	ioc_health "github.com/infraboard/mcube/v2/ioc/apps/health"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	"github.com/infraboard/mcube/v2/ioc/config/http"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

func init() {
	ioc.Api().Registry(&HealthChecker{
		HealthCheck: ioc_health.HealthCheck{
			Path: ioc_health.DEFAUL_HEALTH_PATH,
		},
	})
}

type HealthChecker struct {
	ioc.ObjectImpl
	ioc_health.HealthCheck
	log *zerolog.Logger

	Service healthgrpc.HealthServer `ioc:"autowire=true;namespace=controllers"`
}

func (h *HealthChecker) Name() string {
	return ioc_health.AppName
}

func (h *HealthChecker) Init() error {
	if h.Service == nil {
		h.Service = health.NewServer()
	}

	h.log = log.Sub("health_check")
	h.Registry()
	return nil
}

func (h *HealthChecker) Meta() ioc.ObjectMeta {
	meta := ioc.DefaultObjectMeta()
	meta.CustomPathPrefix = h.Path
	return meta
}

func (h *HealthChecker) Registry() {
	tags := []string{"健康检查"}

	ws := gorestful.ObjectRouter(h)
	ws.Route(ws.GET("/").To(h.HealthHandleFunc).
		Doc("查询服务当前状态").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", ioc_health.HealthCheckResponse{}))

	h.log.Info().Msgf("Get the Health using http://%s%s", http.Get().Addr(), h.Path)
}

func (h *HealthChecker) HealthHandleFunc(r *restful.Request, w *restful.Response) {
	req := ioc_health.NewHealthCheckRequest()
	resp, err := h.Service.Check(
		r.Request.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}

	err = w.WriteAsJson(ioc_health.NewHealth(resp))
	if err != nil {
		h.log.Error().Msgf("send success response error, %s", err)
	}
}
