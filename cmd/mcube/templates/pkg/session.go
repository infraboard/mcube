package pkg

// SessionTemplate todo
const SessionTemplate = `package pkg

import (
	"github.com/infraboard/keyauth/pkg/token"
)

var (
	sess SessionGetter
)

type SessionGetter interface {
	GetToken(requestID string) *token.Token
}

func SetSessionGetter(s SessionGetter) {
	sess = s
}

func S() SessionGetter {
	return sess
}
`
