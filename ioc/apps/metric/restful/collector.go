package metric

import (
	"strconv"
	"time"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/apps/metric"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	ioc_rest "github.com/infraboard/mcube/v2/ioc/config/gorestful"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Api().Registry(&apiStatsCollector{})
}

type apiStatsCollector struct {
	log *zerolog.Logger
	ioc.ObjectImpl

	httpRequestTotal    *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
}

func (h *apiStatsCollector) Init() error {
	h.log = log.Sub(metric.AppName)

	h.httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration",
			Help: "Histogram of the duration of HTTP requests",
			ConstLabels: map[string]string{
				"app": application.Get().AppName,
			},
		},
		[]string{"method", "path"},
	)
	h.httpRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "Total number of HTTP rquests",
			ConstLabels: map[string]string{
				"app": application.Get().AppName,
			},
		},
		[]string{"method", "path", "status_code"},
	)
	// 注册采集器
	prometheus.MustRegister(h)

	// 注册中间件
	ioc_rest.RootRouter().Filter(h.Middleware())
	return nil
}

func (h *apiStatsCollector) Name() string {
	return "api_stats_collector"
}

// func(*Request, *Response, *FilterChain)
func (h *apiStatsCollector) Middleware() restful.FilterFunction {
	return func(r *restful.Request, w *restful.Response, chain *restful.FilterChain) {
		start := time.Now()

		// 处理请求
		chain.ProcessFilter(r, w)
		h.httpRequestDuration.WithLabelValues(r.Request.Method, r.SelectedRoutePath()).Observe(time.Since(start).Seconds())
		h.httpRequestTotal.WithLabelValues(r.Request.Method, r.SelectedRoutePath(), strconv.Itoa(w.StatusCode())).Inc()
	}
}

func (h *apiStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	h.httpRequestDuration.Describe(ch)
	h.httpRequestTotal.Describe(ch)
}

func (h *apiStatsCollector) Collect(ch chan<- prometheus.Metric) {
	h.httpRequestDuration.Collect(ch)
	h.httpRequestTotal.Collect(ch)
}
