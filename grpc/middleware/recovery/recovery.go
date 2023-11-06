package recovery

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RecoveryExplanation 异常消息
const RecoveryExplanation = "Something went wrong"

// Default todo
func Default() *Interceptor {
	return NewInterceptor(NewZeroLogRecoveryHandler())
}

// NewInterceptor todo
func NewInterceptor(h Handler) *Interceptor {
	return &Interceptor{
		h: h,
	}
}

// Interceptor todo
type Interceptor struct {
	h Handler
}

// UnaryServerInterceptor returns a new unary server interceptor for auth.
func (i *Interceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return i.serverIntercept
}

// StreamServerInterceptor todo
func (i *Interceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return i.streamIntercept
}

// Auth impl interface
func (i *Interceptor) serverIntercept(
	ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("%s. Recovering, but please report this.", RecoveryExplanation)
			i.h.Handle(ctx, r)
			// 返回500报错
			err = status.Errorf(codes.Internal, "%v", msg)
			return
		}
	}()

	return handler(ctx, req)
}

func (i Interceptor) streamIntercept(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) (err error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("%s. Recovering, but please report this.", RecoveryExplanation)
			i.h.Handle(nil, r)
			// 返回500报错
			err = status.Errorf(codes.Internal, "%v", msg)
			return
		}
	}()

	err = handler(srv, stream)
	return err
}
