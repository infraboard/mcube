package rest_test

import (
	"context"
	"crypto/tls"
	"testing"

	"github.com/infraboard/mcube/client/negotiator"
	"github.com/infraboard/mcube/client/rest"
	"github.com/infraboard/mcube/http/response"
	"github.com/infraboard/mcube/logger/zap"
	"go.opentelemetry.io/otel"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var (
	ctx = context.Background()
)

func TestClient(t *testing.T) {
	c := rest.NewRESTClient()
	c.EnableTrace()
	c.SetBaseURL("https://www.baidu.com")
	c.SetBearerTokenAuth("UAoVkI07gDGlfARUTToCA8JW")
	var h string
	resp := make(map[string]any)
	err := c.Group("group1").Group("group2").Get("/getpath").Prefix("pre").Suffix("sub").Param("test", "test01").
		Do(context.Background()).
		Header("Server", &h).
		Into(response.NewData(&resp))
	if err != nil {
		t.Fatal(err)
	}

	t.Log(h)
}

func TestTLSClient(t *testing.T) {
	c := rest.NewRESTClient()
	c.SetBaseURL("https://localhost:9443")

	c.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
	resp := make(map[string]any)
	err := c.Post("/mutate--v1-pod").
		Do(ctx).
		ContentType(negotiator.MIME_JSON).
		Into(&resp)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(resp)
}

// 参考样例: https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/net/http/otelhttp/example/client/client.go
func initTracer() error {
	// Create stdout exporter to be able to retrieve
	// the collected spans.
	exporter, err := stdout.New(stdout.WithPrettyPrint())
	if err != nil {
		return err
	}

	// For the demonstration, use sdktrace.AlwaysSample sampler to sample all traces.
	// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSyncer(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("rest client"),
			semconv.ServiceVersionKey.String("v0.0.1"))),
	)
	resource.WithHost()
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return err
}

func init() {
	// 设置日志模式
	zap.DevelopmentSetup()
	// 初始化全局tracer
	err := initTracer()
	if err != nil {
		panic(err)
	}
}
