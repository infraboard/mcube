package httprouter_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/http/router/httprouter"
	"github.com/stretchr/testify/require"
)

func TestSubRouterTestSuit(t *testing.T) {
	suit := newSubRouterTestSuit(t)
	suit.SetUp()
	defer suit.TearDown()

	t.Run("SetLabel", suit.testSetLabel())
	t.Run("AddPublictOK", suit.testAddPublictOK())
	t.Run("ResourceRouterOK", suit.testResourceRouterOK())
}

func newSubRouterTestSuit(t *testing.T) *subRouterTestSuit {
	return &subRouterTestSuit{
		should: require.New(t),
	}
}

type subRouterTestSuit struct {
	root   router.Router
	sub    router.SubRouter
	should *require.Assertions
}

func (s *subRouterTestSuit) SetUp() {
	s.root = httprouter.New()
	s.sub = s.root.SubRouter("/v1")
}

func (a *subRouterTestSuit) TearDown() {

}

func (a *subRouterTestSuit) testSetLabel() func(t *testing.T) {
	return func(t *testing.T) {
		a.sub.SetLabel(router.NewLable("k1", "v1"))
		a.sub.AddPublict("GET", "/index", IndexHandler)

		entries := a.root.GetEndpoints()

		a.should.Equal(entries[0].Labels["k1"], "v1")
	}
}

func (a *subRouterTestSuit) testAddPublictOK() func(t *testing.T) {
	return func(t *testing.T) {
		a.sub.AddPublict("GET", "/", IndexHandler)

		t.Log(a.root.GetEndpoints())

		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/v1/", nil)
		a.root.ServeHTTP(w, req)

		a.should.Equal(w.Code, 200)
	}
}

func (a *subRouterTestSuit) testResourceRouterOK() func(t *testing.T) {
	return func(t *testing.T) {
		a.sub.ResourceRouter("resources").AddPublict("GET", "/", IndexHandler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/resources/", nil)
		a.root.ServeHTTP(w, req)

		t.Log(a.root.GetEndpoints())
		a.should.Equal(w.Code, 200)
	}
}
