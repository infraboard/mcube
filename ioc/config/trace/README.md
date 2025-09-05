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


## 采样策略

TraceIDRatioBased(ratio float64) 确实是一个百分比采样器：参数范围：0.0 到 1.0 之间的浮点数

+ 0.0：不采样任何 trace（0%）
+ 0.5：采样 50% 的请求
+ 1.0：采样所有请求（100%）

工作原理：基于 TraceID 进行哈希计算，确保相同的 TraceID 在不同服务中会有相同的采样决策

推荐的采样策略:
```sh
前端/网关 (采样率: 10%)
    ↓
你的Go服务 (ParentBased + 根采样率: 0%)
    ↓
结果: 只有那10%被网关采样的请求会在你的服务中产生trace
```


## 生产环境注意事项

+ 监控导出器状态：添加健康检查监控导出器是否正常工作
+ 错误处理：妥善处理导出失败的情况，避免影响主业务
+ 性能调优：根据实际流量调整批量处理参数
+ 安全配置：生产环境使用 TLS 和认证
+ 避免 stdout：生产环境不要使用 stdout exporter


## 配置

```toml
[trace]
enable = true
endpoint = "127.0.0.1:4318"
```