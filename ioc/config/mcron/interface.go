package mcron

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/robfig/cron/v3"
)

const (
	APP_NAME = "cron"
	PRIORITY = -199
)

func Get() *cron.Cron {
	return ioc.Default().Get(APP_NAME).(*config).cron
}

func RunAndAddFunc(spec string, cmd func()) (cron.EntryID, error) {
	cmd()

	return Get().AddFunc(spec, cmd)
}
