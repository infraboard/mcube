package gin

import (
	"github.com/gin-gonic/gin"
	h_response "github.com/infraboard/mcube/v2/http/response"
	"github.com/infraboard/mcube/v2/ioc"
	ioc_health "github.com/infraboard/mcube/v2/ioc/apps/health"
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
	return nil
}

func (h *HealthChecker) Meta() ioc.ObjectMeta {
	meta := ioc.DefaultObjectMeta()
	meta.CustomPathPrefix = h.Path
	return meta
}

func (h *HealthChecker) Registry(r gin.IRouter) {
	r.GET("/", h.HealthHandleFunc)

	h.log.Info().Msgf("Get the Health using http://%s%s", http.Get().Addr(), h.Path)
}

func (h *HealthChecker) HealthHandleFunc(c *gin.Context) {
	req := ioc_health.NewHealthCheckRequest()
	resp, err := h.Service.Check(
		c.Request.Context(),
		req,
	)
	if err != nil {
		h_response.Failed(c.Writer, err)
		return
	}

	h_response.Success(c.Writer, ioc_health.NewHealth(resp))
}
