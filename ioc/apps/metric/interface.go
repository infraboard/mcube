package metric

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
	}
}

type Metric struct {
	Enable   bool            `json:"enable" yaml:"enable" toml:"enable" env:"METRIC_ENABLE"`
	Provider METRIC_PROVIDER `toml:"provider" json:"provider" yaml:"provider" env:"METRIC_PROVIDER"`
	Endpoint string          `toml:"endpoint" json:"endpoint" yaml:"endpoint" env:"METRIC_ENDPOINT"`
}
