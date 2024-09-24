package metric

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/apps/metric"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	ioc_gin "github.com/infraboard/mcube/v2/ioc/config/gin"
	"github.com/infraboard/mcube/v2/ioc/config/http"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/prometheus/client_golang/prometheus"
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
	if h.ApiStats.Enable {
		h.AddApiCollector()
	}
	h.Registry()
	return nil
}

func (h *ginHandler) AddApiCollector() {
	collector := metric.NewApiStatsCollector(h.ApiStats, application.Get().AppName)
	// 注册采集器
	prometheus.MustRegister(collector)

	// 注册中间件
	ioc_gin.RootRouter().Use(func(ctx *gin.Context) {
		start := time.Now()

		// 处理请求
		ctx.Next()
		if h.ApiStats.RequestHistogram {
			collector.HttpRequestDurationHistogram.WithLabelValues(ctx.Request.Method, ctx.FullPath()).Observe(time.Since(start).Seconds())
		}
		if h.ApiStats.RequestSummary {
			collector.HttpRequestDurationSummary.WithLabelValues(ctx.Request.Method, ctx.FullPath()).Observe(time.Since(start).Seconds())
		}
		if h.ApiStats.RequestTotal {
			collector.HttpRequestTotal.WithLabelValues(ctx.Request.Method, ctx.FullPath(), strconv.Itoa(ctx.Writer.Status())).Inc()
		}
	})
}

func (h *ginHandler) Name() string {
	return metric.AppName
}

func (h *ginHandler) Version() string {
	return "v1"
}

func (i *ginHandler) Priority() int {
	return 99
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
