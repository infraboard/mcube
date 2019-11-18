package recovery

import (
	"net/http"

	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/logger"
)

const recoveryExplanation = "Something went wrong"

// New returns a new recovery instance
func New() router.Middleware {
	return &recovery{}
}

// NewWithLogger returns a new recovery instance
func NewWithLogger(l logger.RecoveryLogger) router.Middleware {
	return &recovery{l}
}

type recovery struct {
	log logger.RecoveryLogger
}

func (m *recovery) Recover() {
}

// Wrap 实现中间
func (m *recovery) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer m.log.Recover(recoveryExplanation)
		next.ServeHTTP(rw, r)
	})
}
