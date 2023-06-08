package rpc_test

import (
	"context"
	"testing"

	"github.com/infraboard/mcube/logger/zap"

	"{{.PKG}}/clients/rpc"
)

var (
	client *rpc.ClientSet
	ctx    = context.Background()
)

func init() {
	if err := zap.DevelopmentSetup(); err != nil {
		panic(err)
	}

	c, err := rpc.NewClientSetFromEnv()
	if err != nil {
		panic(err)
	}
	client = c
}