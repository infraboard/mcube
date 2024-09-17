package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	AppName = "metric"
)

type METRIC_PROVIDER string

const (
	METRIC_PROVIDER_PROMETHEUS METRIC_PROVIDER = "prometheus"
)

func NewDefaultMetric() *Metric {
	return &Metric{
		Enable:   false,
		Provider: METRIC_PROVIDER_PROMETHEUS,
		Endpoint: "/metrics",
		ApiStatsConfig: ApiStatsConfig{
			ApiRequestHistogram: false,
			ApiRequestSummary:   true,
			ApiRequestTotal:     true,
		},
	}
}

type Metric struct {
	Enable   bool            `json:"enable" yaml:"enable" toml:"enable" env:"ENABLE"`
	Provider METRIC_PROVIDER `toml:"provider" json:"provider" yaml:"provider" env:"PROVIDER"`
	Endpoint string          `toml:"endpoint" json:"endpoint" yaml:"endpoint" env:"ENDPOINT"`
	ApiStatsConfig
}

func NewApiStatsCollector(conf ApiStatsConfig, appName string) *ApiStatsCollector {
	return &ApiStatsCollector{
		ApiStatsConfig: conf,
		HttpRequestDurationHistogram: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "http_request_duration",
				Help: "Histogram of the duration of HTTP requests",
				ConstLabels: map[string]string{
					"app": appName,
				},
			},
			[]string{"method", "path"},
		),
		HttpRequestDurationSummary: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name: "http_request_duration",
				Help: "Histogram of the duration of HTTP requests",
				ConstLabels: map[string]string{
					"app": appName,
				},
				Objectives: map[float64]float64{
					0.5:  0.05,
					0.9:  0.01,
					0.99: 0.001,
				},
			},
			[]string{"method", "path"},
		),
		HttpRequestTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_request_total",
				Help: "Total number of HTTP rquests",
				ConstLabels: map[string]string{
					"app": appName,
				},
			},
			[]string{"method", "path", "status_code"},
		),
	}
}

type ApiStatsConfig struct {
	ApiRequestHistogram bool `toml:"api_request_histogram" json:"api_request_histogram" yaml:"api_request_histogram" env:"API_REQUEST_HISTOGRAM"`
	ApiRequestSummary   bool `toml:"api_request_summary" json:"api_request_summary" yaml:"api_request_summary" env:"API_REQUEST_SUMMARY"`
	ApiRequestTotal     bool `toml:"api_request_total" json:"api_request_total" yaml:"api_request_total" env:"API_REQUEST_TOTAL"`
}

type ApiStatsCollector struct {
	ApiStatsConfig

	HttpRequestTotal             *prometheus.CounterVec
	HttpRequestDurationHistogram *prometheus.HistogramVec
	HttpRequestDurationSummary   *prometheus.SummaryVec
}

func (h *ApiStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	if h.ApiRequestHistogram {
		h.HttpRequestDurationHistogram.Describe(ch)
	}
	if h.ApiRequestSummary {
		h.HttpRequestDurationSummary.Describe(ch)
	}
	if h.ApiRequestTotal {
		h.HttpRequestTotal.Describe(ch)
	}
}

func (h *ApiStatsCollector) Collect(ch chan<- prometheus.Metric) {
	if h.ApiRequestHistogram {
		h.HttpRequestDurationHistogram.Collect(ch)
	}
	if h.ApiRequestSummary {
		h.HttpRequestDurationSummary.Collect(ch)
	}
	if h.ApiRequestTotal {
		h.HttpRequestTotal.Collect(ch)
	}
}
