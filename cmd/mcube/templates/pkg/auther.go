package pkg

// AutherTemplate 模板
const AutherTemplate = `package pkg

import (
	"net/http"

	"github.com/infraboard/mcube/http/router"
)

// NewAuther 内部使用的auther
func NewAuther() router.Auther {
	return &auther{}
}

// internal todo
type auther struct {
}

// 校验身份
func (a *auther) Auth(req *http.Request, entry router.Entry) (
	authInfo interface{}, err error) {

	return nil, nil
}`
