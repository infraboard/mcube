package mcron

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Default().Registry(defaultConfig)
}

var defaultConfig = &config{
	cron: cron.New(cron.WithChain(
		cron.Recover(&LogWrapper{}),
		cron.SkipIfStillRunning(&LogWrapper{}),
	),
		cron.WithLogger(&LogWrapper{}),
	),
}

type config struct {
	cron *cron.Cron
	ioc.ObjectImpl
	log *zerolog.Logger
}

func (c *config) Name() string {
	return APP_NAME
}

func (c *config) Priority() int {
	return PRIORITY
}

func (c *config) Init() error {
	c.log = log.Sub(c.Name())
	c.cron.Start()
	return nil
}
