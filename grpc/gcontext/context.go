package gcontext

import (
	"context"
	"fmt"
	"net/http"

	"github.com/infraboard/mcube/http/request"
	"google.golang.org/grpc/metadata"
)

const (
	// InternalCallTokenHeader todo
	InternalCallTokenHeader = "internal-call-token"
	// ClientIDHeader tood
	ClientIDHeader = "client-id"
	// ClientSecretHeader todo
	ClientSecretHeader = "client-secret"
	// OauthTokenHeader todo
	OauthTokenHeader = "x-oauth-token"
	// RealIPHeader todo
	RealIPHeader = "x-real-ip"
	// UserAgentHeader todo
	UserAgentHeader = "user-agent"
)

// NewGrpcInCtx todo
func NewGrpcInCtx() *GrpcInCtx {
	return &GrpcInCtx{newGrpcCtx(metadata.Pairs())}
}

// GetGrpcInCtx todo
func GetGrpcInCtx(ctx context.Context) (*GrpcInCtx, error) {
	// 获取认证信息
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("ctx is not an grpc incoming context")
	}

	return &GrpcInCtx{newGrpcCtx(md)}, nil
}

// GrpcInCtx todo
type GrpcInCtx struct {
	*grpcCtx
}

// Context todo
func (c *GrpcInCtx) Context() context.Context {
	return metadata.NewIncomingContext(context.Background(), c.md)
}

// SetClientCredentials todo
func (c *GrpcInCtx) SetClientCredentials(clientID, clientSecret string) {
	c.set(ClientIDHeader, clientID)
	c.set(ClientSecretHeader, clientSecret)
}

// GetClientID todo
func (c *GrpcInCtx) GetClientID() string {
	return c.get(ClientIDHeader)
}

// GetClientSecret todo
func (c *GrpcInCtx) GetClientSecret() string {
	return c.get(ClientSecretHeader)
}

// GetAccessToKen todo
func (c *GrpcInCtx) GetAccessToKen() string {
	return c.get(OauthTokenHeader)
}

// NewGrpcOutCtx todo
func NewGrpcOutCtx() *GrpcOutCtx {
	return &GrpcOutCtx{newGrpcCtx(metadata.Pairs())}
}

// GrpcOutCtx todo
type GrpcOutCtx struct {
	*grpcCtx
}

// Context todo
func (c *GrpcOutCtx) Context() context.Context {
	return metadata.NewOutgoingContext(context.Background(), c.md)
}

// SetRemoteIP todo
func (c *GrpcOutCtx) SetRemoteIP(ip string) {
	c.set(RealIPHeader, ip)
}

// SetUserAgent todo
func (c *GrpcOutCtx) SetUserAgent(ua string) {
	c.set(UserAgentHeader, ua)
}

func newGrpcCtx(md metadata.MD) *grpcCtx {
	return &grpcCtx{md: md}
}

// GrpcCtx todo
type grpcCtx struct {
	md metadata.MD
}

// Get todo
func (c *grpcCtx) get(key string) string {
	return c.getWithIndex(key, 0)
}

// Get todo
func (c *grpcCtx) getWithIndex(key string, index int) string {
	if val, ok := c.md[key]; ok {
		if len(val) > index {
			return val[index]
		}
	}

	return ""
}

func (c *grpcCtx) set(key string, values ...string) {
	c.md.Set(key, values...)
}

// SetAccessToken todo
func (c *grpcCtx) SetAccessToken(ak string) {
	c.set(OauthTokenHeader, ak)
}

// NewGrpcOutCtxFromHTTPRequest 从上下文中获取Token
func NewGrpcOutCtxFromHTTPRequest(r *http.Request) (*GrpcOutCtx, error) {
	rc := NewGrpcOutCtx()
	rc.SetAccessToken(r.Header.Get("X-OAUTH-TOKEN"))
	rc.SetRemoteIP(request.GetRemoteIP(r))
	rc.SetUserAgent(r.UserAgent())
	return rc, nil
}
