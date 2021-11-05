package gcontext

import (
	"context"
	"fmt"
	"net/http"

	"github.com/infraboard/mcube/http/request"
	"github.com/rs/xid"
	"google.golang.org/grpc/metadata"
)

const (
	// NamespaceHeader 空间
	NamespaceHeader = "x-rpc-namespace"
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
	// RequestID todo
	RequestID = "x-request-id"
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

// GetNamespace todo
func (c *GrpcInCtx) GetNamespace() string {
	return c.get(NamespaceHeader)
}

// SetAccessToKen 设置ID
func (c *GrpcInCtx) SetAccessToKen(token string) {
	c.set(OauthTokenHeader, token)
}

// GetAccessToKen todo
func (c *GrpcInCtx) GetAccessToKen() string {
	return c.get(OauthTokenHeader)
}

// GetRequestID 获取ID
func (c *GrpcInCtx) GetRequestID() string {
	return c.get(RequestID)
}

// NewGrpcOutCtx todo
func NewGrpcOutCtx() *GrpcOutCtx {
	return &GrpcOutCtx{newGrpcCtx(metadata.Pairs())}
}

// NewGrpcOutCtxFromIn todo
func NewGrpcOutCtxFromIn(in *GrpcInCtx) *GrpcOutCtx {
	return &GrpcOutCtx{newGrpcCtx(in.md)}
}

// GrpcOutCtx todo
type GrpcOutCtx struct {
	*grpcCtx
}

// Context todo
func (c *GrpcOutCtx) Context() context.Context {
	return metadata.NewOutgoingContext(context.Background(), c.md)
}

// SetNamesapce todo
func (c *GrpcOutCtx) SetNamesapce(ns string) {
	c.set(NamespaceHeader, ns)
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

// SetRequestID 设置ID
func (c *grpcCtx) SetRequestID(requestID string) {
	c.set(RequestID, requestID)
}

// SetRemoteIP todo
func (c *grpcCtx) SetRemoteIP(ip string) {
	c.set(RealIPHeader, ip)
}

// SetUserAgent todo
func (c *grpcCtx) SetUserAgent(ua string) {
	c.set(UserAgentHeader, ua)
}

// NewGrpcOutCtxFromHTTPRequest 从上下文中获取Token
func NewGrpcOutCtxFromHTTPRequest(r *http.Request) (*GrpcOutCtx, error) {
	rc := NewGrpcOutCtx()
	rc.SetAccessToken(GetTokenFromHeader(r))
	rc.SetRemoteIP(request.GetRemoteIP(r))
	rc.SetUserAgent(r.UserAgent())
	rc.SetNamesapce(r.Header.Get(NamespaceHeader))
	rc.SetRemoteIP(request.GetRemoteIP(r))

	rid := r.Header.Get(RequestID)
	if rid == "" {
		rid = xid.New().String()
	}
	rc.SetRequestID(rid)

	return rc, nil
}

// NewGrpcOutCtxFromHTTPRequest 从上下文中获取Token
func NewGrpcInCtxFromHTTPRequest(r *http.Request) (*GrpcInCtx, error) {
	rc := NewGrpcInCtx()
	rc.SetAccessToken(GetTokenFromHeader(r))
	rc.SetRemoteIP(request.GetRemoteIP(r))
	rc.SetUserAgent(r.UserAgent())
	rid := r.Header.Get(RequestID)
	if rid == "" {
		rid = xid.New().String()
	}
	rc.SetRequestID(rid)

	return rc, nil
}

func GetTokenFromHeader(r *http.Request) string {
	// 优先从只定义header中读取
	tk := r.Header.Get(OauthTokenHeader)
	if tk != "" {
		return tk
	}

	return r.Header.Get("Authorization")
}

func GetClientCredentialsFromHTTPRequest(r *http.Request) (cid, cs string) {
	cid, cs = r.Header.Get(ClientIDHeader), r.Header.Get(ClientSecretHeader)
	return
}
