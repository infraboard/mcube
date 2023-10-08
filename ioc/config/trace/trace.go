package trace

type Tracer struct {
	TRACE_PROVIDER TRACE_PROVIDER `toml:"provider" json:"provider" yaml:"provider" env:"TRACE_PROVIDER"`
	Endpoint       string         `toml:"endpoint" json:"endpoint" yaml:"endpoint" env:"JAEGER_ENDPOINT"`
}

type TRACE_PROVIDER int

const (
	TRACE_PROVIDER_JAEGER TRACE_PROVIDER = iota
)
