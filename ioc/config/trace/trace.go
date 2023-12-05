package trace

import (
	"github.com/infraboard/mcube/v2/ioc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func init() {
	ioc.Config().Registry(&Tracer{
		Provider: TRACE_PROVIDER_JAEGER,
		Enabled:  true,
	})
}

type Tracer struct {
	Provider TRACE_PROVIDER `toml:"provider" json:"provider" yaml:"provider" env:"TRACE_PROVIDER"`
	Endpoint string         `toml:"endpoint" json:"endpoint" yaml:"endpoint" env:"TRACE_PROVIDER_ENDPOINT"`
	Enabled  bool           `toml:"enabled" json:"enabled" yaml:"enabled" env:"TRACE_ENABLED"`

	ioc.ObjectImpl
}

func (t *Tracer) Name() string {
	return AppName
}

// 优先级最高，要优于日志，日志模块需要依赖trace进行包装
func (i *Tracer) Priority() int {
	return 99
}

func (t *Tracer) Init() error {
	ep := t.Endpoint

	if ep == "" || !t.Enabled {
		return nil
	}

	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(ep)))
	if err != nil {
		return err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.Default()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{}))
	return nil
}
