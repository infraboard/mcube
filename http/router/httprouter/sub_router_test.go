package httprouter_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/http/router/httprouter"
	"github.com/stretchr/testify/require"
)

type subRouterTestSuit struct {
	root router.Router
	sub  router.SubRouter
}

func (s *subRouterTestSuit) SetUp() {
	s.root = httprouter.NewHTTPRouter()
	s.sub = s.root.SubRouter("/v1")
}

func (a *subRouterTestSuit) TearDown() {

}

func TestSubRouterTestSuit(t *testing.T) {
	suit := new(subRouterTestSuit)
	suit.SetUp()
	defer suit.TearDown()

	t.Run("AddPublictOK", testAddPublictOK(suit))
}

func testAddPublictOK(s *subRouterTestSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)

		s.sub.AddPublict("GET", "/", IndexHandler)

		t.Log(s.root.GetEndpoints())
		go http.ListenAndServe("0.0.0.0:8000", s.root)

		<-time.After(time.Second * 10)
		should.NoError(nil)
	}
}
