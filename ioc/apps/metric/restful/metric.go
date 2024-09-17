package metric

import (
	"strconv"
	"time"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/apps/metric"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	ioc_rest "github.com/infraboard/mcube/v2/ioc/config/gorestful"
	"github.com/infraboard/mcube/v2/ioc/config/http"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Api().Registry(&restfulHandler{
		Metric: metric.NewDefaultMetric(),
	})
}

type restfulHandler struct {
	log *zerolog.Logger
	ioc.ObjectImpl

	*metric.Metric
}

func (h *restfulHandler) Init() error {
	h.log = log.Sub(metric.AppName)
	h.Registry()
	if h.ApiStats.Enable {
		h.AddApiCollector()
	}
	return nil
}

func (h *restfulHandler) Name() string {
	return metric.AppName
}

func (h *restfulHandler) Version() string {
	return "v1"
}

func (i *restfulHandler) Priority() int {
	return 1
}

func (h *restfulHandler) Meta() ioc.ObjectMeta {
	meta := ioc.DefaultObjectMeta()
	meta.CustomPathPrefix = h.Endpoint
	return meta
}

func (h *restfulHandler) AddApiCollector() {
	collector := metric.NewApiStatsCollector(h.ApiStats, application.Get().AppName)
	// 注册采集器
	prometheus.MustRegister(collector)

	// 注册中间件
	ioc_rest.RootRouter().Filter(func(r *restful.Request, w *restful.Response, fc *restful.FilterChain) {
		start := time.Now()

		// 处理请求
		fc.ProcessFilter(r, w)
		if h.ApiStats.RequestHistogram {
			collector.HttpRequestDurationHistogram.WithLabelValues(r.Request.Method, r.SelectedRoutePath()).Observe(time.Since(start).Seconds())
		}
		if h.ApiStats.RequestSummary {
			collector.HttpRequestDurationSummary.WithLabelValues(r.Request.Method, r.SelectedRoutePath()).Observe(time.Since(start).Seconds())
		}
		if h.ApiStats.RequestTotal {
			collector.HttpRequestTotal.WithLabelValues(r.Request.Method, r.SelectedRoutePath(), strconv.Itoa(w.StatusCode())).Inc()
		}
	})
}

func (h *restfulHandler) Registry() {
	tags := []string{"健康检查"}
	ws := ioc_rest.ObjectRouter(h)
	ws.Route(ws.
		GET("/").
		To(h.MetricHandleFunc).
		Doc("创建Job").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)

	h.log.Info().Msgf("Get the Metric using http://%s%s", http.Get().Addr(), h.Endpoint)
}

func (h *restfulHandler) MetricHandleFunc(r *restful.Request, w *restful.Response) {
	// 基于标准库 包装了一层
	promhttp.Handler().ServeHTTP(w, r.Request)
}
