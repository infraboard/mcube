package recovery

import (
	"fmt"
	"log"
	"net/http"

	"go.uber.org/zap"

	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/http/response"
	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/logger"
)

const recoveryExplanation = "Something went wrong"

// New returns a new recovery instance
func New() router.Middleware {
	return &recovery{}
}

// NewWithLogger returns a new recovery instance
func NewWithLogger(l logger.WithMetaLogger) router.Middleware {
	return &recovery{
		log: l,
	}
}

type recovery struct {
	log   logger.WithMetaLogger
	debug bool
}

func (m *recovery) Debug(on bool) {
	m.debug = on
}

// Wrap 实现中间
func (m *recovery) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				msg := fmt.Sprintf("%s. Recovering, but please report this.", recoveryExplanation)
				stack := m.stack()

				// 记录Panic日志
				m.logf(msg, r, stack)
				if m.debug {
					msg += stack
				}

				// 返回500报错
				response.Failed(rw, exception.NewInternalServerError(msg))
				return
			}
		}()

		next.ServeHTTP(rw, r)
	})
}

func (m *recovery) stack() string {
	return zap.Stack("stack").String
}

func (m *recovery) logf(msg string, r interface{}, stack string) {
	if m.log != nil {
		m.log.Errorw(msg, logger.NewAny("panic", r), logger.NewAny("stack", stack))
		return
	}

	log.Println(msg, r, stack)
	return
}
