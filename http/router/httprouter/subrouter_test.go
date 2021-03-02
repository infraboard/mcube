package httprouter_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/http/router/httprouter"
	httppb "github.com/infraboard/mcube/pb/http"
	"github.com/stretchr/testify/assert"
)

func TestSubRouterTestSuit(t *testing.T) {
	suit := newSubRouterTestSuit(t)
	suit.SetUp()
	defer suit.TearDown()

	t.Run("SetLabel", suit.testSetLabel())
	t.Run("AddPublictOK", suit.testAddPublictOK())
	t.Run("ResourceRouterOK", suit.testResourceRouterOK())
	t.Run("WithParamsOK", suit.testWithParams())
	t.Run("SetEntryLabel", suit.testSetLabelWithEntry())
}

func newSubRouterTestSuit(t *testing.T) *subRouterTestSuit {
	return &subRouterTestSuit{
		should: assert.New(t),
	}
}

type subRouterTestSuit struct {
	root   router.Router
	sub    router.SubRouter
	should *assert.Assertions
}

func (s *subRouterTestSuit) SetUp() {
	s.root = httprouter.New()
	s.sub = s.root.SubRouter("/v1")
}

func (a *subRouterTestSuit) TearDown() {

}

func (a *subRouterTestSuit) testSetLabel() func(t *testing.T) {
	return func(t *testing.T) {
		a.sub.SetLabel(httppb.NewLable("k1", "v1"))
		a.sub.Handle("GET", "/index", IndexHandler)

		es := a.root.GetEndpoints()

		a.should.Equal(es.GetEntry("/v1/index", "GET").Labels["k1"], "v1")
	}
}

func (a *subRouterTestSuit) testAddPublictOK() func(t *testing.T) {
	return func(t *testing.T) {
		a.sub.Handle("GET", "/", IndexHandler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/", nil)
		a.root.ServeHTTP(w, req)

		a.should.Equal(w.Code, 200)
	}
}

func (a *subRouterTestSuit) testResourceRouterOK() func(t *testing.T) {
	return func(t *testing.T) {
		rr := a.sub.ResourceRouter("resourceName", httppb.NewLable("k1", "v1"))
		rr.BasePath("resources")
		rr.Handle("GET", "/", IndexHandler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/resources/", nil)
		a.root.ServeHTTP(w, req)

		entry := a.root.GetEndpoints().GetEntry("/v1/resources/", "GET")
		a.should.Equal(w.Code, 200)
		a.should.Equal(entry.Labels["k1"], "v1")
	}
}

func (a *subRouterTestSuit) testWithParams() func(t *testing.T) {
	return func(t *testing.T) {
		should := assert.New(t)

		req, _ := http.NewRequest("GET", "/v1/resources/"+urlParam, nil)
		w := httptest.NewRecorder()

		a.sub.Handle("GET", "/resources/:id", WithContextHandler)
		a.root.ServeHTTP(w, req)

		should.Equal(200, w.Code)
	}
}

func (a *subRouterTestSuit) testSetLabelWithEntry() func(t *testing.T) {
	return func(t *testing.T) {
		label := httppb.NewLable("k1", "v1")
		a.sub.Handle("GET", "/index/entry/label", IndexHandler).AddLabel(label)

		es := a.root.GetEndpoints()

		a.should.Equal(es.GetEntry("/v1/index/entry/label", "GET").Labels["k1"], "v1")
	}
}
