package recovery

import (
	"context"
	"log"

	"github.com/infraboard/mcube/logger"
	"go.uber.org/zap"
)

// Handler is a function that recovers from the panic `p` by returning an `error`.
// The context can be used to extract request scoped metadata and context values.
type Handler interface {
	Handle(ctx context.Context, p interface{}) error
}

// NewZapRecoveryHandler todo
func NewZapRecoveryHandler() *ZapRecoveryHandler {
	return &ZapRecoveryHandler{}
}

// ZapRecoveryHandler todo
type ZapRecoveryHandler struct {
	log logger.Logger
}

// SetLogger todo
func (h *ZapRecoveryHandler) SetLogger(l logger.Logger) *ZapRecoveryHandler {
	h.log = l
	return h
}

// Handle todo
func (h *ZapRecoveryHandler) Handle(ctx context.Context, p interface{}) error {
	stack := zap.Stack("stack").String

	if h.log != nil {
		h.log.Errorw(RecoveryExplanation, logger.NewAny("panic", p), logger.NewAny("stack", stack))
		return nil
	}

	log.Println(RecoveryExplanation, p, stack)
	return nil
}
