package rest_test

import (
	"context"
	"testing"

	"github.com/infraboard/mcube/logger/zap"
	"github.com/infraboard/mpaas/client/rest"
	"github.com/infraboard/mpaas/test/tools"
)

var (
	client *rest.ClientSet
	ctx    = context.Background()
)

func init() {
	zap.DevelopmentSetup()
	conf := rest.NewDefaultConfig()
	client = rest.NewClient(conf)
}