package gcontext

import (
	"strconv"

	"github.com/infraboard/mcube/exception"
	"google.golang.org/grpc/metadata"
)

var (
	// Namespace todo
	Namespace = "default"
)

const (
	// ResponseCodeHeader todo
	ResponseCodeHeader = "x-rpc-code"
	// ResponseReasonHeader todo
	ResponseReasonHeader = "x-rpc-reason"
	// ResponseDescHeader todo
	ResponseDescHeader = "x-rpc-desc"
	// ResponseMetaHeader todo
	ResponseMetaHeader = "x-rpc-meta"
	// ResponseDataHeader todo
	ResponseDataHeader = "x-rpc-data"
)

// NewExceptionFromTrailer todo
func NewExceptionFromTrailer(md metadata.MD, err error) exception.APIException {
	ctx := newGrpcCtx(md)
	code, _ := strconv.Atoi(ctx.get(ResponseCodeHeader))
	reason := ctx.get(ResponseReasonHeader)
	message := ctx.get(ResponseDescHeader)
	ctx.get(ResponseMetaHeader)
	ctx.get(ResponseDataHeader)
	if message == "" {
		message = err.Error()
	}
	return exception.NewAPIException(Namespace, code, reason, message)
}
