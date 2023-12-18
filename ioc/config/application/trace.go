package application

import (
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type TRACE_PROVIDER string

const (
	TRACE_PROVIDER_JAEGER TRACE_PROVIDER = "jaeger"
)

func NewDefaultTrace() *Trace {
	return &Trace{
		Provider: TRACE_PROVIDER_JAEGER,
		Endpoint: "http://localhost:14268/api/traces",
		Enable:   false,
	}
}

type Trace struct {
	Enable   bool           `json:"enable" yaml:"enable" toml:"enable" env:"ENABLE"`
	Provider TRACE_PROVIDER `toml:"provider" json:"provider" yaml:"provider" env:"TRACE_PROVIDER"`
	Endpoint string         `toml:"endpoint" json:"endpoint" yaml:"endpoint" env:"TRACE_PROVIDER_ENDPOINT"`
}

func (t *Trace) SetDefaultEnv() {
	sn := os.Getenv("OTEL_SERVICE_NAME")
	if sn == "" {
		os.Setenv("OTEL_SERVICE_NAME", App().AppName)
	}
}

func (t *Trace) Parse() error {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(t.Endpoint)))
	if err != nil {
		return err
	}

	t.SetDefaultEnv()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.Default()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{}))
	return nil
}
