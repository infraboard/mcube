package mock

import (
	"net/http"
	"strings"

	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/http/router"
	httppb "github.com/infraboard/mcube/pb/http"
)

var (
	// MockTestToken 用于内部mock测试
	MockTestToken = "justformocktest"
)

// NewMockAuther mock
func NewMockAuther() router.Auther {
	return &mockAuther{}
}

type mockAuther struct{}

func (m *mockAuther) Auth(r *http.Request, entry httppb.Entry) (authInfo interface{}, err error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, exception.NewUnauthorized("Authorization missed in header")
	}

	headerSlice := strings.Split(authHeader, " ")
	if len(headerSlice) != 2 {
		return nil, exception.NewUnauthorized("Authorization header value is not validated, must be: {token_type} {token}")
	}

	access := headerSlice[1]

	if access != MockTestToken {
		return nil, exception.NewUnauthorized("permission deny")
	}
	return access, nil
}

func (m *mockAuther) ResponseHook(http.ResponseWriter, *http.Request, httppb.Entry) {

}
