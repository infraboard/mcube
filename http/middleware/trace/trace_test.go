package trace_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"

	"github.com/infraboard/mcube/http/middleware/trace"
	"github.com/infraboard/mcube/http/router/httprouter"
	"github.com/infraboard/mcube/ioc/config/logger"
	"github.com/infraboard/mcube/logger/zap"
)

var (
	tracer = mocktracer.New()
	peer   = trace.PeerInfo{
		Service:  "mock-trace",
		Address:  "0.0.0.0",
		Port:     8080,
		Hostname: "MacBook-Pro-2.local",
	}
)

func Test_Trace(t *testing.T) {
	should := assert.New(t)

	tr := trace.New(tracer, peer)
	tr.Debug(true)
	tr.SetLogger(logger.Sub("Trace"))

	r := httprouter.New()
	r.Use(tr)
	r.Handle("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	// 准备请求
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	spans := tracer.FinishedSpans()
	if should.Equal(1, len(spans)) {
		target := spans[0]
		tags := target.Tags()
		should.Equal(trace.DefaultComponentName, tags[string(ext.Component)])
		should.Equal("GET", tags[string(ext.HTTPMethod)])
		should.Equal(uint16(200), tags[string(ext.HTTPStatusCode)])
		should.Equal(peer.Service, tags[string(ext.PeerService)])
		should.Equal("/", target.OperationName)
		should.Equal(req.Header.Get(trace.RequestIDHeaderKey), w.Header().Get(trace.RequestIDHeaderKey))
		t.Log(target.Tags())
	}
}

// func Test_TraceWithClient(t *testing.T) {
// 	cfg := &config.Configuration{
// 		Sampler: &config.SamplerConfig{
// 			Type:  "const",
// 			Param: 1.0,
// 		},
// 		Reporter: &config.ReporterConfig{
// 			LocalAgentHostPort: "127.0.0.1:6831",
// 			LogSpans:           true,
// 		},
// 	}

// 	jaegerS, closerS, _ := cfg.New("mock-server")
// 	defer closerS.Close()

// 	opentracing.SetGlobalTracer(jaegerS)

// 	tr := trace.New(jaegerS, peer)
// 	tr.Debug(true)
// 	tr.SetLogger(logger.Sub("Trace"))

// 	r := httprouter.New()
// 	r.Use(tr)
// 	r.Handle("GET", "/", func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte("Hello World"))
// 	})
// 	svr := httptest.NewServer(r)

// 	// 准备请求
// 	jaegerC, closerC, _ := cfg.New("mock-clientxxxx")
// 	defer closerC.Close()

// 	req, _ := http.NewRequest("GET", svr.URL, nil)
// 	req.Header.Add("X-Request-Id", "request id")
// 	span := jaegerC.StartSpan(req.URL.String())
// 	defer span.Finish()
// 	req = req.WithContext(opentracing.ContextWithSpan(req.Context(), span))
// 	req, ht := nethttp.TraceRequest(jaegerC, req)
// 	defer ht.Finish()

// 	client := &http.Client{Transport: &nethttp.Transport{}}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	_ = resp.Body.Close()

// 	svr.Close()
// }

func init() {
	zap.DevelopmentSetup()
}
