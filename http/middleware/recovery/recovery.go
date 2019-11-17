package recovery

import (
	"net/http"

	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/logger"
)

const recoveryExplanation = "Something went wrong"

// New returns a new recovery instance
func New(l logger.RecoveryLogger) router.Middleware {
	return &recoveryMiddleWare{l}
}

type recoveryMiddleWare struct {
	logger.RecoveryLogger
}

// Wrap 实现中间
func (m *recoveryMiddleWare) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer m.Recover(recoveryExplanation)
		next.ServeHTTP(rw, r)
	})
}
