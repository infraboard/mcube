package httprouter_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/http/context"
	"github.com/infraboard/mcube/http/mock"
	"github.com/infraboard/mcube/http/response"
	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/http/router/httprouter"
	"github.com/stretchr/testify/assert"
)

const (
	urlParam = "101"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	response.Success(w, "ok")
	return
}

func WithContextHandler(w http.ResponseWriter, r *http.Request) {
	rctx := context.GetContext(r)
	id := rctx.PS.ByName("id")

	if id == urlParam {
		response.Success(w, "ok")
		return
	}

	response.Failed(w, exception.NewBadRequest("failed"))
	return
}

func TestBase(t *testing.T) {
	should := assert.New(t)

	r := httprouter.New()

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	r.AddPublict("GET", "/", IndexHandler)
	r.ServeHTTP(w, req)
	should.Equal(w.Code, 200)
}

func TestWithAutherFailed(t *testing.T) {
	should := assert.New(t)

	r := httprouter.New()
	// 设置mock
	r.SetAuther(mock.NewMockAuther())

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "bearer invalid_token")
	w := httptest.NewRecorder()

	r.AddProtected("GET", "/", IndexHandler)
	r.ServeHTTP(w, req)

	should.Equal(401, w.Code)
}

func TestWithAutherOK(t *testing.T) {
	should := assert.New(t)

	r := httprouter.New()
	// 设置mock
	r.SetAuther(mock.NewMockAuther())

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "bearer "+mock.MockTestToken)
	w := httptest.NewRecorder()

	r.AddProtected("GET", "/", IndexHandler)
	r.ServeHTTP(w, req)

	should.Equal(200, w.Code)
}

func TestWithParamsOK(t *testing.T) {
	should := assert.New(t)

	r := httprouter.New()

	req, _ := http.NewRequest("GET", "/"+urlParam, nil)
	w := httptest.NewRecorder()

	r.AddProtected("GET", "/:id", WithContextHandler)
	r.ServeHTTP(w, req)

	should.Equal(200, w.Code)
}

func TestWithParamsFailed(t *testing.T) {
	should := assert.New(t)

	r := httprouter.New()

	req, _ := http.NewRequest("GET", "/"+"urlParam", nil)
	w := httptest.NewRecorder()

	r.AddProtected("GET", "/:id", WithContextHandler)
	r.ServeHTTP(w, req)

	should.Equal(400, w.Code)
}

func TestSetLabel(t *testing.T) {
	should := assert.New(t)

	r := httprouter.New()
	r.SetLabel(router.NewLable("k1", "v1"))
	r.AddProtected("GET", "/:id", WithContextHandler)

	es := r.GetEndpoints()

	should.Equal(es.GetEntry("/:id", "GET").Labels["k1"], "v1")
}

func TestAPIRootOK(t *testing.T) {
	should := assert.New(t)

	r := httprouter.New()

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	r.AddProtected("GET", "/test", IndexHandler)
	r.EnableAPIRoot()
	r.ServeHTTP(w, req)

	es := router.NewEntrySet()
	if should.NoError(response.GetDataFromBody(w.Result().Body, es)) {
		should.NotNil(es.GetEntry("/test", "GET"))
	}
}

func TestAPIRootOrderOK(t *testing.T) {
	should := assert.New(t)

	r := httprouter.New()

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	r.AddProtected("GET", "/test1", IndexHandler)
	r.AddProtected("GET", "/test2", IndexHandler)
	r.AddPublict("GET", "/test3", IndexHandler)
	r.EnableAPIRoot()
	r.ServeHTTP(w, req)

	es := router.NewEntrySet()
	if should.NoError(response.GetDataFromBody(w.Result().Body, es)) {
		should.Equal("/test1", es.Items[0].Path)
		should.Equal("/test2", es.Items[1].Path)
		should.Equal(true, es.Items[1].Protected)
		should.Equal("/test3", es.Items[2].Path)
		should.Equal(false, es.Items[2].Protected)

	}
}
