package trace

import "github.com/infraboard/mcube/v2/ioc"

const (
	AppName = "trace"
)

type TRACE_PROVIDER string

const (
	// 标准otlp http协议
	TRACE_PROVIDER_OTLP TRACE_PROVIDER = "otlp"
	// 开发环境使用
	TRACE_PROVIDER_STDOUT TRACE_PROVIDER = "stdout"
)

func Get() *Trace {
	obj := ioc.Config().Get(AppName)
	if obj == nil {
		return defaultConfig
	}
	return obj.(*Trace)
}
