package conf

// ToolsTempate todo
const ToolsTempate = `package conf

import (
	"context"

	"github.com/infraboard/keyauth/pkg/token"
	"github.com/infraboard/mcube/grpc/gcontext"
)

func GetTokenFromGrpcInCtx(ctx context.Context) (*token.Token, error) {
	kc, err := C().Keyauth.Client()
	if err != nil {
		return nil, err
	}
	in, err := gcontext.GetGrpcInCtx(ctx)
	if err != nil {
		return nil, err
	}
	req := token.NewDescribeTokenRequestWithAccessToken(in.GetAccessToKen())
	out := gcontext.NewGrpcOutCtx()
	out.SetAccessToken(in.GetAccessToKen())
	return kc.Token().DescribeToken(out.Context(), req)
}
`
