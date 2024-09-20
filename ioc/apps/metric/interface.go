package metric

import (
	"time"

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
		Enable:   true,
		Provider: METRIC_PROVIDER_PROMETHEUS,
		Endpoint: "/metrics",
		ApiStats: ApiStatsConfig{
			Enable:                     true,
			AppConstLabelKey:           "app",
			RequestHistogram:           false,
			RequestHistogramName:       "http_request_duration",
			RequestHistogramBucket:     prometheus.DefBuckets,
			RequestSummary:             true,
			RequestSummaryName:         "request_latency_percentile",
			RequestSummaryObjective:    []float64{0.5, 0.9, 0.99},
			RequestSummaryMaxAgeSecond: 10 * 60,
			RequestTotal:               true,
			RequestTotalName:           "http_request_total",
		},
	}
}

type Metric struct {
	Enable   bool            `json:"enable" yaml:"enable" toml:"enable" env:"ENABLE"`
	Provider METRIC_PROVIDER `toml:"provider" json:"provider" yaml:"provider" env:"PROVIDER"`
	Endpoint string          `toml:"endpoint" json:"endpoint" yaml:"endpoint" env:"ENDPOINT"`
	ApiStats ApiStatsConfig  `toml:"api_stats" json:"api_stats" yaml:"api_stats" env:"API_STATS"`
}

func NewApiStatsCollector(conf ApiStatsConfig, appName string) *ApiStatsCollector {
	return &ApiStatsCollector{
		ApiStatsConfig: conf,
		HttpRequestDurationHistogram: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: conf.RequestHistogramName,
				Help: "Histogram of the duration of HTTP requests",
				ConstLabels: map[string]string{
					conf.AppConstLabelKey: appName,
				},
				Buckets: conf.RequestHistogramBucket,
			},
			[]string{"method", "path"},
		),
		HttpRequestDurationSummary: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name: conf.RequestSummaryName,
				Help: "Histogram of the duration of HTTP requests",
				ConstLabels: map[string]string{
					conf.AppConstLabelKey: appName,
				},
				Objectives: conf.Objectives(),
				MaxAge:     time.Duration(conf.RequestSummaryMaxAgeSecond) * time.Second,
			},
			[]string{"method", "path"},
		),
		HttpRequestTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: conf.RequestTotalName,
				Help: "Total number of HTTP rquests",
				ConstLabels: map[string]string{
					conf.AppConstLabelKey: appName,
				},
			},
			[]string{"method", "path", "status_code"},
		),
	}
}

type ApiStatsConfig struct {
	Enable                     bool      `toml:"enable" json:"enable" yaml:"enable" env:"ENABLE"`
	AppConstLabelKey           string    `toml:"app_const_label_key" json:"app_const_label_key" yaml:"app_const_label_key" env:"APP_CONST_LABEL_KEY"`
	RequestHistogram           bool      `toml:"request_histogram" json:"request_histogram" yaml:"request_histogram" env:"REQUEST_HISTOGRAM"`
	RequestHistogramName       string    `toml:"request_histogram_name" json:"request_histogram_name" yaml:"request_histogram_name" env:"REQUEST_HISTOGRAM_NAME"`
	RequestHistogramBucket     []float64 `toml:"request_histogram_bucket" json:"request_histogram_bucket" yaml:"request_histogram_bucket" env:"REQUEST_HISTOGRAM_BUCKET" envSeparator:","`
	RequestSummary             bool      `toml:"request_summary" json:"request_summary" yaml:"request_summary" env:"REQUEST_SUMMARY"`
	RequestSummaryName         string    `toml:"request_summary_name" json:"request_summary_name" yaml:"request_summary_name" env:"REQUEST_SUMMARY_NAME"`
	RequestSummaryObjective    []float64 `toml:"request_summary_objective" json:"request_summary_objective" yaml:"request_summary_objective" env:"REQUEST_SUMMARY_OBJECTIVE" envSeparator:","`
	RequestSummaryMaxAgeSecond int64     `toml:"request_summary_max_age_second" json:"request_summary_max_age_second" yaml:"request_summary_max_age_second" env:"REQUEST_SUMMARY_MAX_AGE_SECOND"`

	RequestTotal     bool   `toml:"request_total" json:"request_total" yaml:"request_total" env:"REQUEST_TOTAL"`
	RequestTotalName string `toml:"request_total_name" json:"request_total_name" yaml:"request_total_name" env:"REQUEST_TOTAL_NAME"`
}

func (c ApiStatsConfig) Objectives() map[float64]float64 {
	objectives := map[float64]float64{}
	for _, v := range c.RequestSummaryObjective {
		objectives[v] = (1 - v) * 0.1
	}
	return objectives
}

type ApiStatsCollector struct {
	ApiStatsConfig
	HttpRequestTotal             *prometheus.CounterVec
	HttpRequestDurationHistogram *prometheus.HistogramVec
	HttpRequestDurationSummary   *prometheus.SummaryVec
}

func (h *ApiStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	if h.RequestHistogram {
		h.HttpRequestDurationHistogram.Describe(ch)
	}
	if h.RequestSummary {
		h.HttpRequestDurationSummary.Describe(ch)
	}
	if h.RequestTotal {
		h.HttpRequestTotal.Describe(ch)
	}
}

func (h *ApiStatsCollector) Collect(ch chan<- prometheus.Metric) {
	if h.RequestHistogram {
		h.HttpRequestDurationHistogram.Collect(ch)
	}
	if h.RequestSummary {
		h.HttpRequestDurationSummary.Collect(ch)
	}
	if h.RequestTotal {
		h.HttpRequestTotal.Collect(ch)
	}
}
