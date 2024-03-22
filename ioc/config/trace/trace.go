package trace

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

func init() {
	ioc.Config().Registry(defaultConfig)
}

var defaultConfig = &Trace{
	Provider: TRACE_PROVIDER_OTLP,
	Endpoint: "localhost:4137",
	Enable:   false,
}

type Trace struct {
	ioc.ObjectImpl

	Enable   bool           `json:"enable" yaml:"enable" toml:"enable" env:"TRACE_ENABLE"`
	Provider TRACE_PROVIDER `toml:"provider" json:"provider" yaml:"provider" env:"TRACE_PROVIDER"`
	Endpoint string         `toml:"endpoint" json:"endpoint" yaml:"endpoint" env:"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT"`

	tp *trace.TracerProvider
}

// 优先初始化, 以供后面的组件使用
func (i *Trace) Priority() int {
	return 99
}

func (i *Trace) Name() string {
	return AppName
}

// otlp go sdk 使用方法: https://opentelemetry.io/docs/languages/go/exporters/
// jaeger 端口说明: https://www.jaegertracing.io/docs/1.55/getting-started/#all-in-one
func (t *Trace) Init() error {
	// 创建一个OTLP exporter
	exporter, err := otlptracegrpc.New(
		context.Background(),
		otlptracegrpc.WithEndpoint(t.Endpoint),
	)
	if err != nil {
		return err
	}

	t.tp = trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
		trace.WithResource(resource.Default()),
	)
	otel.SetTracerProvider(t.tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{}))
	return nil
}

func (t *Trace) Close(ctx context.Context) error {
	return t.tp.Shutdown(ctx)
}
