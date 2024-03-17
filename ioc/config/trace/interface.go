package trace

import "github.com/infraboard/mcube/v2/ioc"

const (
	AppName = "trace"
)

type TRACE_PROVIDER string

const (
	TRACE_PROVIDER_OTLP TRACE_PROVIDER = "otlp"
)

func Get() *Trace {
	obj := ioc.Config().Get(AppName)
	if obj == nil {
		return defaultConfig
	}
	return obj.(*Trace)
}
