package trace

import (
	"log"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/rs/xid"
	"github.com/rs/zerolog"

	"github.com/infraboard/mcube/http/response"
)

const (
	// RequestIDHeaderKey 补充的RequestID Header Key
	RequestIDHeaderKey = "X-Request-Id"
	// DefaultComponentName 默认组件名称
	DefaultComponentName = "mcube.server.trace"
)

// New 实例
func New(tr opentracing.Tracer, peer PeerInfo) *Tracer {
	opentracing.SetGlobalTracer(tr)
	return &Tracer{
		tr:    tr,
		debug: true,
		peer:  peer,
	}
}

// PeerInfo 端的信息
type PeerInfo struct {
	Service  string
	Address  string
	Port     uint16
	Hostname string
}

// Tracer 链路追踪器
type Tracer struct {
	tr    opentracing.Tracer
	log   *zerolog.Logger
	debug bool
	peer  PeerInfo
}

// Debug 是否开启debug
func (h *Tracer) Debug(on bool) {
	h.debug = on
}

// SetLogger 设置Logger
func (h *Tracer) SetLogger(l *zerolog.Logger) {
	h.log = l
}

// Handler todo
func (h *Tracer) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var span opentracing.Span

		// Attempt to join a trace by getting trace context from the headers.
		carrier := opentracing.HTTPHeadersCarrier(req.Header)
		ctx, err := h.tr.Extract(opentracing.HTTPHeaders, carrier)
		if err != nil {
			span = opentracing.StartSpan(req.URL.String(), ext.RPCServerOption(nil))
		} else {
			span = opentracing.StartSpan(req.URL.String(), ext.RPCServerOption(ctx))
		}

		// 添加Request ID
		requestID := req.Header.Get(RequestIDHeaderKey)
		if requestID == "" {
			requestID = xid.New().String()
			req.Header.Set(RequestIDHeaderKey, requestID)
		}

		span.SetTag("request.id", requestID)
		if err := span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier); err != nil {
			ext.Error.Set(span, true)
		} else {
			req = req.WithContext(opentracing.ContextWithSpan(req.Context(), span))
		}

		// 补充基本信息
		if span != nil {
			ext.HTTPMethod.Set(span, req.Method)
			ext.Component.Set(span, DefaultComponentName)
			ext.PeerService.Set(span, h.peer.Service)
			ext.PeerAddress.Set(span, h.peer.Address)
			ext.PeerHostname.Set(span, h.peer.Hostname)
			ext.PeerPort.Set(span, h.peer.Port)
		}

		// 结束时补充status code
		defer func() {
			if span != nil {
				res := rw.(response.Response)
				res.Header().Set(RequestIDHeaderKey, requestID)
				ext.HTTPStatusCode.Set(span, uint16(res.Status()))
				span.Finish()
			}
			h.logf("%s", span)
		}()

		next.ServeHTTP(rw, req)
	})
}

func (h *Tracer) logf(format string, args ...interface{}) {
	if !h.debug {
		return
	}

	if h.log != nil {
		h.log.Debug().Msgf(format, args...)
		return
	}

	log.Printf(format+"\n", args...)
}
