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
	tr.SetLogger(zap.L().Named("Trace"))

	r := httprouter.New()
	r.Use(tr)
	r.Handle("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	// 准备请求
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("User-Agent", "mcube/0.3")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	spans := tracer.FinishedSpans()
	if should.Equal(1, len(spans)) {
		target := spans[0]
		tags := target.Tags()
		should.Equal("mcube/trace", tags[string(ext.Component)])
		should.Equal("GET", tags[string(ext.HTTPMethod)])
		should.Equal(uint16(200), tags[string(ext.HTTPStatusCode)])
		should.Equal(peer.Service, tags[string(ext.PeerService)])
		should.Equal("/", target.OperationName)
	}
}

func init() {
	zap.DevelopmentSetup()
}
