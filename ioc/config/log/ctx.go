package log

import (
	"context"
	"strings"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

func FromCtx(ctx context.Context, names ...string) *zerolog.Logger {
	l := zerolog.Ctx(ctx)
	if Get().IsMyLogger(l) {
		return l
	}

	if len(names) == 0 {
		l = L()
	} else {
		l = Sub(strings.Join(names, "."))
	}

	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() && spanCtx.HasTraceID() {
		traceID := spanCtx.TraceID().String()
		var tl = l.Hook(zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, message string) {
			e.Fields(map[string]any{
				Get().TraceFiled: traceID,
			})
		}))
		return &tl
	}
	return l
}
