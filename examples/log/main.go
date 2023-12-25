package main

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/log"
)

func main() {
	ioc.DevelopmentSetup()

	gLogger := log.L()
	gLogger.Debug().Msg("this is global logger debug msg")

	subLogger := log.Sub("app1")
	subLogger.Debug().Msg("this is app1 sub logger debug msg")

	ctx := context.Background()
	traceLogger := log.T("app1").Trace(ctx)
	traceLogger.Debug().Msg("this is app1 trace logger debug msg")
}
