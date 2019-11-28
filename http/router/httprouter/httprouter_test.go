package httprouter_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/infraboard/mcube/http/auth/mock"
	"github.com/infraboard/mcube/http/response"
	"github.com/infraboard/mcube/http/router/httprouter"
	"github.com/stretchr/testify/require"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	response.Success(w, http.StatusOK, "ok")
	return
}

func TestBase(t *testing.T) {
	should := require.New(t)

	r := httprouter.New()

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	r.AddProtected("GET", "/", IndexHandler)
	r.ServeHTTP(w, req)
	should.Equal(w.Code, 200)
}

func TestWithAutherFailed(t *testing.T) {
	should := require.New(t)

	r := httprouter.New()
	// 设置mock
	r.SetAuther(mock.NewMockAuther())

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "bearer invalid_token")
	w := httptest.NewRecorder()

	r.AddProtected("GET", "/", IndexHandler)
	r.ServeHTTP(w, req)

	should.Equal(w.Code, 403)
}

func TestWithAutherOK(t *testing.T) {
	should := require.New(t)

	r := httprouter.New()
	// 设置mock
	r.SetAuther(mock.NewMockAuther())

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "bearer "+mock.MockTestToken)
	w := httptest.NewRecorder()

	r.AddProtected("GET", "/", IndexHandler)
	r.ServeHTTP(w, req)

	should.Equal(w.Code, 200)
}
