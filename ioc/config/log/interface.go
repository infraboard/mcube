package log

import (
	"context"
	"strings"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/rs/zerolog"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	AppName = "log"
)

const (
	SUB_LOGGER_KEY = "component"
)

func L() *zerolog.Logger {
	return Get().root
}

func Sub(name string) *zerolog.Logger {
	return Get().Logger(name)
}

func Ctx(ctx context.Context, names ...string) *zerolog.Logger {
	var l *zerolog.Logger
	if len(names) == 0 {
		l = L()
	} else {
		l = Sub(strings.Join(names, "."))
	}

	traceId := oteltrace.SpanFromContext(ctx).SpanContext().TraceID().String()
	if traceId != "" {
		var tl = l.Hook(zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, message string) {
			e.Fields(map[string]any{
				Get().TraceFiled: traceId,
			})
		}))
		return &tl
	}

	return l
}

func Get() *Config {
	obj := ioc.Config().Get(AppName)
	if obj == nil {
		return defaultConfig
	}
	return ioc.Config().Get(AppName).(*Config)
}
