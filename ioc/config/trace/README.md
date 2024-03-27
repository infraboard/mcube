# Trace

## Jaeger
To try out the OTLP exporter, since v1.35.0 you can run Jaeger as an OTLP endpoint and for trace visualization in a Docker container:

```sh
docker run -d --name jaeger \
  -e COLLECTOR_OTLP_ENABLED=true \
  -p 16686:16686 \
  -p 4317:4317 \
  -p 4318:4318 \
  jaegertracing/all-in-one:latest
```

+ otlp go sdk 使用方法: https://opentelemetry.io/docs/languages/go/exporters/
+ jaeger 端口说明: https://www.jaegertracing.io/docs/1.55/getting-started/#all-in-one