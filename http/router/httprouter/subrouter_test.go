package httprouter_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/http/router/httprouter"
	"github.com/stretchr/testify/require"
)

func TestSubRouterTestSuit(t *testing.T) {
	suit := newSubRouterTestSuit(t)
	suit.SetUp()
	defer suit.TearDown()

	t.Run("AddPublictOK", suit.testAddPublictOK())
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
	s.root = httprouter.NewHTTPRouter()
	s.sub = s.root.SubRouter("/v1")
}

func (a *subRouterTestSuit) TearDown() {

}

func (a *subRouterTestSuit) testAddPublictOK() func(t *testing.T) {
	return func(t *testing.T) {
		a.sub.AddPublict("GET", "/", IndexHandler)

		t.Log(a.root.GetEndpoints())
		go http.ListenAndServe("0.0.0.0:8000", a.root)

		<-time.After(time.Second * 2)
		a.should.NoError(nil)
	}
}
