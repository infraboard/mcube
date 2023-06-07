package rpc_test

import (
	"context"
	"os"
	"testing"

	"github.com/infraboard/mcenter/apps/endpoint"
	"github.com/infraboard/mcenter/apps/instance"
	"github.com/infraboard/mcenter/apps/token"
	"github.com/infraboard/mcenter/client/rpc"
	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/ioc/health"
)

var (
	c   *rpc.ClientSet
	ctx = context.Background()
)


func init() {
	err := rpc.LoadClientFromEnv()
	if err != nil {
		panic(err)
	}
	c = rpc.C()
}