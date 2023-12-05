package accesslog_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/infraboard/mcube/v2/http/middleware/accesslog"
	"github.com/infraboard/mcube/v2/http/router"
	"github.com/infraboard/mcube/v2/http/router/httprouter"
	"github.com/stretchr/testify/require"
)

func TestAccessLogTestSuit(t *testing.T) {
	suit := new(accessLogTestSuit)
	suit.SetUp()
	defer suit.TearDown()

	t.Run("LogOK", suit.Test_LoggerOK())
	t.Run("URLEncodedStringOK", suit.Test_LoggerURLEncodedString())
	t.Run("CustomFormatOK", suit.Test_LoggerCustomFormat())
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "hello world")
}

type accessLogTestSuit struct {
	router router.Router
	lm     *accesslog.Logger
}

func (s *accessLogTestSuit) SetUp() {
	s.router = httprouter.New()

	lm := accesslog.New()
	s.lm = lm
	s.router.Use(lm)
	s.router.Handle("GET", "/", indexHandler)
}

func (s *accessLogTestSuit) TearDown() {

}

func (s *accessLogTestSuit) Test_LoggerOK() func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)
		recorder := httptest.NewRecorder()

		req, err := http.NewRequest("GET", "http://localhost:3000/", nil)
		should.NoError(err)

		s.router.ServeHTTP(recorder, req)

		logMsg := ""
		s.lm.SetDebugFunc(func(msg string) {
			logMsg = msg
		})

		should.Contains(logMsg, "localhost:3000 | GET /")
	}
}

func (s *accessLogTestSuit) Test_LoggerURLEncodedString() func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)

		recorder := httptest.NewRecorder()

		// Test reserved characters - !*'();:@&=+$,/?%#[]
		req, err := http.NewRequest("GET", "http://localhost:3000/%21%2A%27%28%29%3B%3A%40%26%3D%2B%24%2C%2F%3F%25%23%5B%5D", nil)
		should.NoError(err)

		s.router.ServeHTTP(recorder, req)

		logMsg := ""
		s.lm.SetDebugFunc(func(msg string) {
			logMsg = msg
		})

		should.Equal(recorder.Code, http.StatusNotFound)
		should.Contains(logMsg, "/!*'();:@&=+$,/?%#[]")
	}
}

func (s *accessLogTestSuit) Test_LoggerCustomFormat() func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)

		recorder := httptest.NewRecorder()

		s.lm.SetFormat("{{.Request.URL.Query.Get \"foo\"}} {{.Request.UserAgent}} - {{.Status}}")
		userAgent := "mcube-test"
		req, err := http.NewRequest("GET", "http://localhost:3000/?foo=bar", nil)
		should.NoError(err)
		req.Header.Set("User-Agent", userAgent)

		s.router.ServeHTTP(recorder, req)
		logMsg := ""
		s.lm.SetDebugFunc(func(msg string) {
			logMsg = msg
		})

		should.Equal(logMsg, "bar "+userAgent+" - 200")
	}
}
