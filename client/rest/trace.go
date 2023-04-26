package rest

import (
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// 开启后一定要配置全局Tracer
func (c *RESTClient) EnableTrace() *RESTClient {
	c.client.Transport = otelhttp.NewTransport(c.transport)
	c.provider = otel.GetTracerProvider()
	c.propagators = otel.GetTextMapPropagator()

	c.tr = c.provider.Tracer(
		"github.com/infraboard/mcube/client/rest",
		oteltrace.WithInstrumentationVersion("semver:0.0.1"),
	)
	return c
}

// 关闭Trace
func (c *RESTClient) DisableTrace() *RESTClient {
	c.client.Transport = c.transport
	c.tr = nil
	return c
}
