package accesslog_test

import (
	"testing"

	"github.com/infraboard/mcube/http/middleware/accesslog"
	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/http/router/httprouter"
)

type accessLogTestSuit struct {
	r router.Router
}

func (a *accessLogTestSuit) SetUp() {
	a.r = httprouter.NewHTTPRouter()
}

func (a *accessLogTestSuit) TearDown() {

}

func (a *accessLogTestSuit) Test_LoggerOK() func(t *testing.T) {
	return func(t *testing.T) {

		lm := accesslog.NewLogger()
		a.r.Use(lm)

		// var buff bytes.Buffer
		// recorder := httptest.NewRecorder()

		// l.ALogger = log.New(&buff, "[negroni] ", 0)
		// n := New()
		// // replace log for testing
		// a
		// n.UseHandler(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// 	rw.WriteHeader(http.StatusNotFound)
		// }))

		// req, err := http.NewRequest("GET", "http://localhost:3000/foobar", nil)
		// if err != nil {
		// 	t.Error(err)
		// }

		// n.ServeHTTP(recorder, req)
		// expect(t, recorder.Code, http.StatusNotFound)
		// refute(t, len(buff.String()), 0)
	}
}

func (a *accessLogTestSuit) Test_LoggerURLEncodedString() func(t *testing.T) {
	return func(t *testing.T) {
	}
}

func (a *accessLogTestSuit) Test_LoggerCustomFormat() func(t *testing.T) {
	return func(t *testing.T) {
	}
}

// func Test_Logger(t *testing.T) {
// }

// func Test_LoggerURLEncodedString(t *testing.T) {
// 	var buff bytes.Buffer
// 	recorder := httptest.NewRecorder()

// 	l := NewLogger()
// 	l.ALogger = log.New(&buff, "[negroni] ", 0)
// 	l.SetFormat("{{.Path}}")

// 	n := New()
// 	// replace log for testing
// 	n.Use(l)
// 	n.UseHandler(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
// 		rw.WriteHeader(http.StatusOK)
// 	}))

// 	// Test reserved characters - !*'();:@&=+$,/?%#[]
// 	req, err := http.NewRequest("GET", "http://localhost:3000/%21%2A%27%28%29%3B%3A%40%26%3D%2B%24%2C%2F%3F%25%23%5B%5D", nil)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	n.ServeHTTP(recorder, req)
// 	expect(t, recorder.Code, http.StatusOK)
// 	expect(t, strings.TrimSpace(buff.String()), "[negroni] /!*'();:@&=+$,/?%#[]")
// 	refute(t, len(buff.String()), 0)
// }

// func Test_LoggerCustomFormat(t *testing.T) {
// 	var buff bytes.Buffer
// 	recorder := httptest.NewRecorder()

// 	l := NewLogger()
// 	l.ALogger = log.New(&buff, "[negroni] ", 0)
// 	l.SetFormat("{{.Request.URL.Query.Get \"foo\"}} {{.Request.UserAgent}} - {{.Status}}")

// 	n := New()
// 	n.Use(l)
// 	n.UseHandler(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
// 		rw.Write([]byte("OK"))
// 	}))

// 	userAgent := "Negroni-Test"
// 	req, err := http.NewRequest("GET", "http://localhost:3000/foobar?foo=bar", nil)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	req.Header.Set("User-Agent", userAgent)

// 	n.ServeHTTP(recorder, req)
// 	expect(t, strings.TrimSpace(buff.String()), "[negroni] bar "+userAgent+" - 200")
// }
