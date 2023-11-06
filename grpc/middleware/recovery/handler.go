package recovery

import (
	"context"
	"runtime/debug"

	"github.com/infraboard/mcube/ioc/config/logger"
	"github.com/rs/zerolog"
)

// Handler is a function that recovers from the panic `p` by returning an `error`.
// The context can be used to extract request scoped metadata and context values.
type Handler interface {
	Handle(ctx context.Context, p interface{}) error
}

// NewZeroLogRecoveryHandler todo
func NewZeroLogRecoveryHandler() *ZeroLogRecoveryHandler {
	return &ZeroLogRecoveryHandler{
		log: logger.Sub("grpc"),
	}
}

// ZeroLogRecoveryHandler todo
type ZeroLogRecoveryHandler struct {
	log *zerolog.Logger
}

// SetLogger todo
func (h *ZeroLogRecoveryHandler) SetLogger(l *zerolog.Logger) *ZeroLogRecoveryHandler {
	h.log = l
	return h
}

// Handle todo
func (h *ZeroLogRecoveryHandler) Handle(ctx context.Context, p interface{}) error {
	// 捕获 panic
	defer func() {
		if err := recover(); err != nil {
			h.log.Error().Stack().Msgf("Panic occurred: %v\n%s", err, debug.Stack())
		}
	}()
	return nil
}
