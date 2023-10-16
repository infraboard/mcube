package cors_test

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/infraboard/mcube/http/middleware/cors"
)

var testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("bar"))
})

var allHeaders = []string{
	"Vary",
	"Access-Control-Allow-Origin",
	"Access-Control-Allow-Methods",
	"Access-Control-Allow-Headers",
	"Access-Control-Allow-Credentials",
	"Access-Control-Max-Age",
	"Access-Control-Expose-Headers",
}

func assertHeaders(t *testing.T, resHeaders http.Header, expHeaders map[string]string) {
	for _, name := range allHeaders {
		got := strings.Join(resHeaders[name], ", ")
		want := expHeaders[name]
		if got != want {
			t.Errorf("Response header %q = %q, want %q", name, got, want)
		}
	}
}

func TestSpec(t *testing.T) {
	cases := []struct {
		name       string
		options    cors.Options
		method     string
		reqHeaders map[string]string
		resHeaders map[string]string
	}{
		{
			"NoConfig",
			cors.Options{
				// Intentionally left blank.
			},
			"GET",
			map[string]string{},
			map[string]string{
				"Vary": "Origin",
			},
		},
		{
			"MatchAllOrigin",
			cors.Options{
				AllowedOrigins: []string{"*"},
			},
			"GET",
			map[string]string{
				"Origin": "http://foobar.com",
			},
			map[string]string{
				"Vary":                        "Origin",
				"Access-Control-Allow-Origin": "*",
			},
		},
		{
			"MatchAllOriginWithCredentials",
			cors.Options{
				AllowedOrigins:   []string{"*"},
				AllowCredentials: true,
			},
			"GET",
			map[string]string{
				"Origin": "http://foobar.com",
			},
			map[string]string{
				"Vary":                             "Origin",
				"Access-Control-Allow-Origin":      "*",
				"Access-Control-Allow-Credentials": "true",
			},
		},
		{
			"AllowedOrigin",
			cors.Options{
				AllowedOrigins: []string{"http://foobar.com"},
			},
			"GET",
			map[string]string{
				"Origin": "http://foobar.com",
			},
			map[string]string{
				"Vary":                        "Origin",
				"Access-Control-Allow-Origin": "http://foobar.com",
			},
		},
		{
			"WildcardOrigin",
			cors.Options{
				AllowedOrigins: []string{"http://*.bar.com"},
			},
			"GET",
			map[string]string{
				"Origin": "http://foo.bar.com",
			},
			map[string]string{
				"Vary":                        "Origin",
				"Access-Control-Allow-Origin": "http://foo.bar.com",
			},
		},
		{
			"DisallowedOrigin",
			cors.Options{
				AllowedOrigins: []string{"http://foobar.com"},
			},
			"GET",
			map[string]string{
				"Origin": "http://barbaz.com",
			},
			map[string]string{
				"Vary": "Origin",
			},
		},
		{
			"DisallowedWildcardOrigin",
			cors.Options{
				AllowedOrigins: []string{"http://*.bar.com"},
			},
			"GET",
			map[string]string{
				"Origin": "http://foo.baz.com",
			},
			map[string]string{
				"Vary": "Origin",
			},
		},
		{
			"AllowedOriginFuncMatch",
			cors.Options{
				AllowOriginFunc: func(o string) bool {
					return regexp.MustCompile("^http://foo").MatchString(o)
				},
			},
			"GET",
			map[string]string{
				"Origin": "http://foobar.com",
			},
			map[string]string{
				"Vary":                        "Origin",
				"Access-Control-Allow-Origin": "http://foobar.com",
			},
		},
		{
			"AllowOriginRequestFuncMatch",
			cors.Options{
				AllowOriginRequestFunc: func(r *http.Request, o string) bool {
					return regexp.MustCompile("^http://foo").MatchString(o) && r.Header.Get("Authorization") == "secret"
				},
			},
			"GET",
			map[string]string{
				"Origin":        "http://foobar.com",
				"Authorization": "secret",
			},
			map[string]string{
				"Vary":                        "Origin",
				"Access-Control-Allow-Origin": "http://foobar.com",
			},
		},
		{
			"AllowOriginRequestFuncNotMatch",
			cors.Options{
				AllowOriginRequestFunc: func(r *http.Request, o string) bool {
					return regexp.MustCompile("^http://foo").MatchString(o) && r.Header.Get("Authorization") == "secret"
				},
			},
			"GET",
			map[string]string{
				"Origin":        "http://foobar.com",
				"Authorization": "not-secret",
			},
			map[string]string{
				"Vary": "Origin",
			},
		},
		{
			"MaxAge",
			cors.Options{
				AllowedOrigins: []string{"http://example.com/"},
				AllowedMethods: []string{"GET"},
				MaxAge:         10,
			},
			"OPTIONS",
			map[string]string{
				"Origin":                        "http://example.com/",
				"Access-Control-Request-Method": "GET",
			},
			map[string]string{
				"Vary":                         "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
				"Access-Control-Allow-Origin":  "http://example.com/",
				"Access-Control-Allow-Methods": "GET",
				"Access-Control-Max-Age":       "10",
			},
		},
		{
			"AllowedMethod",
			cors.Options{
				AllowedOrigins: []string{"http://foobar.com"},
				AllowedMethods: []string{"PUT", "DELETE"},
			},
			"OPTIONS",
			map[string]string{
				"Origin":                        "http://foobar.com",
				"Access-Control-Request-Method": "PUT",
			},
			map[string]string{
				"Vary":                         "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
				"Access-Control-Allow-Origin":  "http://foobar.com",
				"Access-Control-Allow-Methods": "PUT",
			},
		},
		{
			"DisallowedMethod",
			cors.Options{
				AllowedOrigins: []string{"http://foobar.com"},
				AllowedMethods: []string{"PUT", "DELETE"},
			},
			"OPTIONS",
			map[string]string{
				"Origin":                        "http://foobar.com",
				"Access-Control-Request-Method": "PATCH",
			},
			map[string]string{
				"Vary": "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
			},
		},
		{
			"AllowedHeaders",
			cors.Options{
				AllowedOrigins: []string{"http://foobar.com"},
				AllowedHeaders: []string{"X-Header-1", "x-header-2"},
			},
			"OPTIONS",
			map[string]string{
				"Origin":                         "http://foobar.com",
				"Access-Control-Request-Method":  "GET",
				"Access-Control-Request-Headers": "X-Header-2, X-HEADER-1",
			},
			map[string]string{
				"Vary":                         "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
				"Access-Control-Allow-Origin":  "http://foobar.com",
				"Access-Control-Allow-Methods": "GET",
				"Access-Control-Allow-Headers": "X-Header-2, X-Header-1",
			},
		},
		{
			"DefaultAllowedHeaders",
			cors.Options{
				AllowedOrigins: []string{"http://foobar.com"},
				AllowedHeaders: []string{},
			},
			"OPTIONS",
			map[string]string{
				"Origin":                         "http://foobar.com",
				"Access-Control-Request-Method":  "GET",
				"Access-Control-Request-Headers": "X-Requested-With",
			},
			map[string]string{
				"Vary":                         "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
				"Access-Control-Allow-Origin":  "http://foobar.com",
				"Access-Control-Allow-Methods": "GET",
				"Access-Control-Allow-Headers": "X-Requested-With",
			},
		},
		{
			"AllowedWildcardHeader",
			cors.Options{
				AllowedOrigins: []string{"http://foobar.com"},
				AllowedHeaders: []string{"*"},
			},
			"OPTIONS",
			map[string]string{
				"Origin":                         "http://foobar.com",
				"Access-Control-Request-Method":  "GET",
				"Access-Control-Request-Headers": "X-Header-2, X-HEADER-1",
			},
			map[string]string{
				"Vary":                         "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
				"Access-Control-Allow-Origin":  "http://foobar.com",
				"Access-Control-Allow-Methods": "GET",
				"Access-Control-Allow-Headers": "X-Header-2, X-Header-1",
			},
		},
		{
			"DisallowedHeader",
			cors.Options{
				AllowedOrigins: []string{"http://foobar.com"},
				AllowedHeaders: []string{"X-Header-1", "x-header-2"},
			},
			"OPTIONS",
			map[string]string{
				"Origin":                         "http://foobar.com",
				"Access-Control-Request-Method":  "GET",
				"Access-Control-Request-Headers": "X-Header-3, X-Header-1",
			},
			map[string]string{
				"Vary": "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
			},
		},
		{
			"OriginHeader",
			cors.Options{
				AllowedOrigins: []string{"http://foobar.com"},
			},
			"OPTIONS",
			map[string]string{
				"Origin":                         "http://foobar.com",
				"Access-Control-Request-Method":  "GET",
				"Access-Control-Request-Headers": "origin",
			},
			map[string]string{
				"Vary":                         "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
				"Access-Control-Allow-Origin":  "http://foobar.com",
				"Access-Control-Allow-Methods": "GET",
				"Access-Control-Allow-Headers": "Origin",
			},
		},
		{
			"ExposedHeader",
			cors.Options{
				AllowedOrigins: []string{"http://foobar.com"},
				ExposedHeaders: []string{"X-Header-1", "x-header-2"},
			},
			"GET",
			map[string]string{
				"Origin": "http://foobar.com",
			},
			map[string]string{
				"Vary":                          "Origin",
				"Access-Control-Allow-Origin":   "http://foobar.com",
				"Access-Control-Expose-Headers": "X-Header-1, X-Header-2",
			},
		},
		{
			"AllowedCredentials",
			cors.Options{
				AllowedOrigins:   []string{"http://foobar.com"},
				AllowCredentials: true,
			},
			"OPTIONS",
			map[string]string{
				"Origin":                        "http://foobar.com",
				"Access-Control-Request-Method": "GET",
			},
			map[string]string{
				"Vary":                             "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
				"Access-Control-Allow-Origin":      "http://foobar.com",
				"Access-Control-Allow-Methods":     "GET",
				"Access-Control-Allow-Credentials": "true",
			},
		},
		{
			"OptionPassthrough",
			cors.Options{
				OptionsPassthrough: true,
			},
			"OPTIONS",
			map[string]string{
				"Origin":                        "http://foobar.com",
				"Access-Control-Request-Method": "GET",
			},
			map[string]string{
				"Vary":                         "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET",
			},
		},
		{
			"NonPreflightOptions",
			cors.Options{
				AllowedOrigins: []string{"http://foobar.com"},
			},
			"OPTIONS",
			map[string]string{
				"Origin": "http://foobar.com",
			},
			map[string]string{
				"Vary":                        "Origin",
				"Access-Control-Allow-Origin": "http://foobar.com",
			},
		},
	}
	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			s := cors.New(tc.options)

			req, _ := http.NewRequest(tc.method, "http://example.com/foo", nil)
			for name, value := range tc.reqHeaders {
				req.Header.Add(name, value)
			}

			t.Run("Handler", func(t *testing.T) {
				res := httptest.NewRecorder()
				s.Handler(testHandler).ServeHTTP(res, req)
				assertHeaders(t, res.Header(), tc.resHeaders)
			})
		})
	}
}

func TestDefault(t *testing.T) {
	s := cors.Default()
	if !s.AllowedOriginsAll() {
		t.Error("c.allowedOriginsAll should be true when Default")
	}
	if s.AllowedHeaders() == nil {
		t.Error("c.allowedHeaders should be nil when Default")
	}
	if s.AllowedMethods() == nil {
		t.Error("c.allowedMethods should be nil when Default")
	}
}

func TestHandlePreflightInvalidOriginAbortion(t *testing.T) {
	s := cors.New(cors.Options{
		AllowedOrigins: []string{"http://foo.com"},
	})
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "http://example.com/foo", nil)
	req.Header.Add("Origin", "http://example.com/")

	s.HandlePreflight(res, req)

	assertHeaders(t, res.Header(), map[string]string{
		"Vary": "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
	})
}

func TestHandlePreflightNoOptionsAbortion(t *testing.T) {
	s := cors.New(cors.Options{
		// Intentionally left blank.
	})
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)

	s.HandlePreflight(res, req)

	assertHeaders(t, res.Header(), map[string]string{})
}

func TestHandleActualRequestInvalidOriginAbortion(t *testing.T) {
	s := cors.New(cors.Options{
		AllowedOrigins: []string{"http://foo.com"},
	})
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	req.Header.Add("Origin", "http://example.com/")

	s.HandleActualRequest(res, req)

	assertHeaders(t, res.Header(), map[string]string{
		"Vary": "Origin",
	})
}

func TestHandleActualRequestInvalidMethodAbortion(t *testing.T) {
	s := cors.New(cors.Options{
		AllowedMethods:   []string{"POST"},
		AllowCredentials: true,
	})
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	req.Header.Add("Origin", "http://example.com/")

	s.HandleActualRequest(res, req)

	assertHeaders(t, res.Header(), map[string]string{
		"Vary": "Origin",
	})
}

func TestIsMethodAllowedReturnsFalseWithNoMethods(t *testing.T) {
	s := cors.New(cors.Options{
		AllowedMethods: []string{},
		// Intentionally left blank.
	})
	if s.IsMethodAllowed("") {
		t.Error("IsMethodAllowed should return false when c.allowedMethods is nil.")
	}
}

func TestIsMethodAllowedReturnsTrueWithOptions(t *testing.T) {
	s := cors.New(cors.Options{
		// Intentionally left blank.
	})
	if !s.IsMethodAllowed("OPTIONS") {
		t.Error("IsMethodAllowed should return true when c.allowedMethods is nil.")
	}
}
