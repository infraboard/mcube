package tracer

import (
	"context"

	"{{.PKG}}/conf"
	"{{.PKG}}/version"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var Tracer oteltrace.Tracer

func InitTracer() error {
	ep := conf.C().Jaeger.Endpoint

	if ep == "" {
		return nil
	}

	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(ep)))
	if err != nil {
		return err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(newResource()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{}))
	Tracer = otel.GetTracerProvider().Tracer(version.ServiceName, oteltrace.WithInstrumentationVersion("0.1"))
	return nil
}

// newResource returns a resource describing this application.
func newResource() *resource.Resource {
	// 因为需要手动设置ServiceName和版本, 不使用default的自动探测
	// resource.Default()

	r, _ := resource.New(
		context.Background(),
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(version.ServiceName),
			semconv.ServiceVersionKey.String(version.Short()),
		),
	)
	return r
}
