package recovery

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/http/response"
	"github.com/infraboard/mcube/v2/http/router"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/rs/zerolog"
)

const recoveryExplanation = "Something went wrong"

// New returns a new recovery instance
func New() router.Middleware {
	return &recovery{
		logger: log.Sub("recovery"),
	}
}

type recovery struct {
	logger *zerolog.Logger
}

// Wrap 实现中间
func (m *recovery) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				msg := fmt.Sprintf("%s. Recovering, but please report this.", recoveryExplanation)

				// 记录Panic日志
				m.logger.Error().Stack().Msgf("Panic occurred: %v\n%s", msg, debug.Stack())

				// 返回500报错
				response.Failed(rw, exception.NewInternalServerError(msg))
				return
			}
		}()

		next.ServeHTTP(rw, r)
	})
}
