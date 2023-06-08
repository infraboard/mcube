package exception

import (
	"context"

	"github.com/infraboard/mcube/exception"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// grpc server端 的异常只支持 code 与 description, 为了能完整把异常传递给下游调用方, 把异常放到了grpc response header中
// 客户端如果发现有这个key 说明该异常是我们定义的业务异常对象(API Exception), 需要还原。
func NewUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return (&UnaryClientInterceptor{}).UnaryClientInterceptor
}

type UnaryClientInterceptor struct {
}

func (e *UnaryClientInterceptor) UnaryClientInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) error {

	var trailer metadata.MD
	opts = append(opts, grpc.Trailer(&trailer))
	err := invoker(ctx, method, req, reply, cc, opts...)
	t := trailer.Get(exception.TRAILER_ERROR_JSON_KEY)
	if len(t) > 0 {
		err = exception.NewAPIExceptionFromString(t[0])
	}
	return err
}
