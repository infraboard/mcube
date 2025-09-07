package jsonrpc_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/ioc/config/jsonrpc"
)

func TestCall(t *testing.T) {
	jsonrpc.RegisterService(&UserService{})
}

// 用户服务
type UserService struct{}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GetUserRequest struct {
	UserID int `json:"userId"`
}

// RPC 方法：自动注册为 getUser
func (s *UserService) RPCGetUser(ctx context.Context, req *GetUserRequest) (*User, error) {
	// 直接使用解析好的请求对象
	if req.UserID <= 0 {
		return nil, exception.NewApiException(1000, "Invalid user ID")
	}

	return &User{
		ID:   req.UserID,
		Name: fmt.Sprintf("User%d", req.UserID),
	}, nil
}
