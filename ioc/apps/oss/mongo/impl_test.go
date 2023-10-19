package mongo_test

import (
	"context"

	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/ioc/apps/oss"
	"github.com/infraboard/mflow/test/tools"
)

var (
	impl oss.Service
	ctx  = context.Background()
)

func init() {
	tools.DevelopmentSetup()
	impl = ioc.GetController(oss.AppName).(oss.Service)
}
