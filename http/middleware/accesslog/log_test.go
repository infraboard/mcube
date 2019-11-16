package accesslog_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/infraboard/mcube/http/middleware/accesslog"
	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/http/router/httprouter"
	"github.com/infraboard/mcube/logger/mock"
	"github.com/stretchr/testify/require"
)

func TestAccessLogTestSuit(t *testing.T) {
	suit := new(accessLogTestSuit)
	suit.SetUp()
	defer suit.TearDown()

	t.Run("PutOK", suit.Test_LoggerOK())
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}

type accessLogTestSuit struct {
	router   router.Router
	mkLogger *mock.StandardLogger
}

func (s *accessLogTestSuit) SetUp() {
	s.router = httprouter.NewHTTPRouter()
	s.mkLogger = mock.NewStandardLogger()
}

func (s *accessLogTestSuit) TearDown() {

}

func (s *accessLogTestSuit) Test_LoggerOK() func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)

		lm := accesslog.NewLogger()
		lm.SetLogger(s.mkLogger)
		s.router.Use(lm)
		s.router.AddPublict("GET", "/", IndexHandler)

		recorder := httptest.NewRecorder()

		req, err := http.NewRequest("GET", "http://localhost:3000/", nil)
		if err != nil {
			t.Error(err)
		}
		s.router.ServeHTTP(recorder, req)

		msg, err := ioutil.ReadAll(s.mkLogger.Buffer)
		should.NoError(err)
		should.Contains(string(msg), "localhost:3000 | GET /")
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
