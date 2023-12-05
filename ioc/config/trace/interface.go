package trace

import "github.com/infraboard/mcube/v2/ioc"

const (
	AppName = "trace"
)

func C() *Config {
	return ioc.Config().Get(AppName).(*Config)
}
