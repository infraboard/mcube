package main

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/logger"
)

func main() {
	ioc.DevelopmentSetup()

	gLogger := logger.L()
	gLogger.Debug().Msg("this is global logger debug msg")

	subLogger := logger.Sub("app1")
	subLogger.Debug().Msg("this is app1 sub logger debug msg")

	ctx := context.Background()
	traceLogger := logger.T("app1").Trace(ctx)
	traceLogger.Debug().Msg("this is app1 trace logger debug msg")
}
