package metric

import (
	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/apps/metric"
	ioc_gin "github.com/infraboard/mcube/v2/ioc/config/gin"
	"github.com/infraboard/mcube/v2/ioc/config/http"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Api().Registry(&ginHandler{
		Metric: metric.NewDefaultMetric(),
	})
}

type ginHandler struct {
	log *zerolog.Logger
	ioc.ObjectImpl

	*metric.Metric
}

func (h *ginHandler) Init() error {
	h.log = log.Sub(metric.AppName)

	h.Registry()
	return nil
}

func (h *ginHandler) Name() string {
	return metric.AppName
}

func (h *ginHandler) Version() string {
	return "v1"
}

func (h *ginHandler) Meta() ioc.ObjectMeta {
	meta := ioc.DefaultObjectMeta()
	meta.CustomPathPrefix = h.Endpoint
	return meta
}

func (h *ginHandler) Registry() {
	r := ioc_gin.ObjectRouter(h)
	r.GET("/", func(ctx *gin.Context) {
		// 基于标准库 包装了一层
		promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
	})

	h.log.Info().Msgf("Get the Metric using http://%s%s", http.Get().Addr(), h.Endpoint)
}
