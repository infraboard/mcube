package mongo_test

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/apps/oss"
)

var (
	impl oss.Service
	ctx  = context.Background()
)

func init() {
	ioc.DevelopmentSetup()
	impl = ioc.GetController(oss.AppName).(oss.Service)
}
