package trace

import (
	"context"
	"io"
	"time"

	"github.com/infraboard/mcube/v2/ioc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

func init() {
	ioc.Config().Registry(defaultConfig)
}

var defaultConfig = &Trace{
	Provider:           TRACE_PROVIDER_STDOUT,
	Endpoint:           "127.0.0.1:4318",
	Insecure:           true,
	Enable:             false,
	TraceIDRatio:       0,
	BatchTimeout:       5 * time.Second,
	MaxExportBatchSize: 512,
	MaxQueueSize:       2048,
}

type Trace struct {
	ioc.ObjectImpl

	Enable   bool           `json:"enable" yaml:"enable" toml:"enable" env:"ENABLE"`
	Provider TRACE_PROVIDER `toml:"provider" json:"provider" yaml:"provider" env:"PROVIDER"`
	Endpoint string         `toml:"endpoint" json:"endpoint" yaml:"endpoint" env:"ENDPOINT"`
	Insecure bool           `toml:"insecure" json:"insecure" yaml:"insecure" env:"INSECURE"`

	// 采样率配置
	TraceIDRatio float64 `toml:"trace_id_ratio" json:"trace_id_ratio" yaml:"trace_id_ratio" env:"TRACE_ID_RATIO"`
	// 上报配置
	BatchTimeout       time.Duration `toml:"batch_timeout" json:"batch_timeout" yaml:"batch_timeout" env:"BATCH_TIMEOUT"`
	MaxExportBatchSize int           `toml:"max_export_batch_size" json:"max_export_batch_size" yaml:"max_export_batch_size" env:"MAX_EXPORT_BATCH_SIZE"`
	MaxQueueSize       int           `toml:"max_queue_size" json:"max_queue_size" yaml:"max_queue_size" env:"MAX_QUEUE_SIZE"`

	tp *trace.TracerProvider
}

// 优先初始化, 以供后面的组件使用
func (i *Trace) Priority() int {
	return 998
}

func (i *Trace) Name() string {
	return AppName
}

func (i *Trace) options() (opts []otlptracehttp.Option) {
	opts = append(opts,
		otlptracehttp.WithEndpoint(i.Endpoint),
		otlptracehttp.WithTimeout(10*time.Second),
		otlptracehttp.WithRetry(otlptracehttp.RetryConfig{
			Enabled:         true,
			InitialInterval: 2 * time.Second,
			MaxInterval:     30 * time.Second,
			MaxElapsedTime:  60 * time.Second,
		}),
	)
	if i.Insecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	return
}

// otlp go sdk 使用方法: https://opentelemetry.io/docs/languages/go/exporters/
// jaeger 端口说明: https://www.jaegertracing.io/docs/1.55/getting-started/#all-in-one
func (t *Trace) Init() error {
	if t.Enable {
		var exporter trace.SpanExporter
		var err error
		// 创建一个OTLP exporter
		switch t.Provider {
		case TRACE_PROVIDER_OTLP:
			exporter, err = otlptracehttp.New(
				context.Background(),
				t.options()...,
			)
		default:
			exporter, err = stdouttrace.New(
				stdouttrace.WithWriter(io.Discard),
				stdouttrace.WithPrettyPrint(),
			)
		}
		if err != nil {
			return err
		}

		// 2. 配置采样器 - 这是实现你需求的关键！
		// ParentBased 采样器会尊重父级 trace 的采样决策。
		// TraceIDRatioBased(0) 参数表示：对于没有父级 trace 的根 span，采样率为 0（即不采样）。
		// 对于有父级 trace 的请求，如果父级采样了，则本服务也采样；如果父级没采样，则本服务也不采样。
		sampler := trace.ParentBased(trace.TraceIDRatioBased(t.TraceIDRatio))

		t.tp = trace.NewTracerProvider(
			trace.WithSampler(sampler),
			trace.WithBatcher(exporter,
				trace.WithBatchTimeout(t.BatchTimeout),
				trace.WithMaxExportBatchSize(t.MaxExportBatchSize),
				trace.WithMaxQueueSize(t.MaxQueueSize)),
			trace.WithResource(resource.Default()),
		)
		otel.SetTracerProvider(t.tp)

		// 5. 设置全局的 Propagator（必须设置，才能从 HTTP 头中提取 trace context）
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, // W3C Trace Context
			propagation.Baggage{},      // W3C Baggage
		))
	}
	return nil
}

func (t *Trace) Close(ctx context.Context) {
	if t.tp != nil {
		t.tp.Shutdown(ctx)
		return
	}
}
