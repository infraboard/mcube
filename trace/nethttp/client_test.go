package nethttp_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"

	"github.com/infraboard/mcube/trace/nethttp"
)

func makeRequest(t *testing.T, url string, options ...nethttp.ClientOption) []*mocktracer.MockSpan {
	should := assert.New(t)

	tr := &mocktracer.MockTracer{}
	span := tr.StartSpan("toplevel")
	client := &http.Client{Transport: &nethttp.Transport{}}
	req, err := http.NewRequest("GET", url, nil)
	should.NoError(err)

	req = req.WithContext(opentracing.ContextWithSpan(req.Context(), span))
	req, ht := nethttp.TraceRequest(tr, req, options...)
	resp, err := client.Do(req)
	should.NoError(err)

	_ = resp.Body.Close()
	ht.Finish()
	span.Finish()

	return tr.FinishedSpans()
}

func TestClientTrace(t *testing.T) {
	should := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ok", http.StatusTemporaryRedirect)
	})
	mux.HandleFunc("/fail", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "failure", http.StatusInternalServerError)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	helloWorldObserver := func(s opentracing.Span, r *http.Request) {
		s.SetTag("hello", "world")
	}

	tests := []struct {
		url          string
		num          int
		opts         []nethttp.ClientOption
		opName       string
		expectedTags map[string]interface{}
	}{
		{url: "/ok", num: 3, opts: nil, opName: "HTTP Client"},
		{url: "/redirect", num: 4, opts: []nethttp.ClientOption{nethttp.OperationName("client-span")}, opName: "client-span"},
		{url: "/fail", num: 3, opts: nil, opName: "HTTP Client", expectedTags: makeTags(string(ext.Error), true)},
		{url: "/ok", num: 3, opts: []nethttp.ClientOption{nethttp.ClientSpanObserver(helloWorldObserver)}, opName: "HTTP Client", expectedTags: makeTags("hello", "world")},
	}

	for _, tt := range tests {
		t.Log(tt.opName)
		spans := makeRequest(t, srv.URL+tt.url, tt.opts...)
		should.Equal(tt.num, len(spans))
		var rootSpan *mocktracer.MockSpan
		for _, span := range spans {
			if span.ParentID == 0 {
				rootSpan = span
				break
			}
		}
		should.NotNil(rootSpan)

		foundClientSpan := false
		for _, span := range spans {
			if span.ParentID == rootSpan.SpanContext.SpanID {
				foundClientSpan = true
				should.Equal(tt.opName, span.OperationName)
			}
			if span.OperationName == "HTTP GET" {
				logs := span.Logs()
				should.LessOrEqual(6, len(logs))
				key := logs[0].Fields[0].Key
				should.Equal("event", key)
				v := logs[0].Fields[0].ValueString
				should.Equal("GetConn", v)

				for k, expected := range tt.expectedTags {
					result := span.Tag(k)
					should.Equal(expected, result)
				}
			}
		}
		should.Equal(true, foundClientSpan)
	}
}

func TestTracerFromRequest(t *testing.T) {
	should := assert.New(t)
	req, err := http.NewRequest("GET", "foobar", nil)
	should.NoError(err)

	ht := nethttp.TracerFromRequest(req)
	should.Nil(ht)

	tr := &mocktracer.MockTracer{}
	req, expected := nethttp.TraceRequest(tr, req)

	ht = nethttp.TracerFromRequest(req)
	should.Equal(expected, ht)
}

var (
	injectTests = []struct {
		name                     string
		expectContextPropagation bool
		opts                     []nethttp.ClientOption
	}{
		{name: "Default", expectContextPropagation: true, opts: nil},
		{name: "True", expectContextPropagation: true, opts: []nethttp.ClientOption{nethttp.InjectSpanContext(true)}},
		{name: "False", expectContextPropagation: false, opts: []nethttp.ClientOption{nethttp.InjectSpanContext(false)}},
	}
)

func TestInjectSpanContext(t *testing.T) {
	for _, tt := range injectTests {
		t.Run(tt.name, func(t *testing.T) {
			var handlerCalled bool
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true
				srvTr := mocktracer.New()
				ctx, err := srvTr.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))

				if err != nil && tt.expectContextPropagation {
					t.Fatal(err)
				}

				if tt.expectContextPropagation {
					if err != nil || ctx == nil {
						t.Fatal("expected propagation but unable to extract")
					}
				} else {
					// Expect "opentracing: SpanContext not found in Extract carrier" when not injected
					// Can't check ctx directly, because it gets set to emptyContext
					if err == nil {
						t.Fatal("unexpected propagation")
					}
				}
			}))

			tr := mocktracer.New()
			span := tr.StartSpan("root")

			req, err := http.NewRequest("GET", srv.URL, nil)
			if err != nil {
				t.Fatal(err)
			}
			req = req.WithContext(opentracing.ContextWithSpan(req.Context(), span))

			req, ht := nethttp.TraceRequest(tr, req, tt.opts...)

			client := &http.Client{Transport: &nethttp.Transport{}}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			_ = resp.Body.Close()

			ht.Finish()
			span.Finish()

			srv.Close()

			if !handlerCalled {
				t.Fatal("server handler never called")
			}
		})
	}
}

func makeTags(keyVals ...interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(keyVals)/2)
	for i := 0; i < len(keyVals)-1; i += 2 {
		key := keyVals[i].(string)
		result[key] = keyVals[i+1]
	}
	return result
}
