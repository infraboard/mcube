package trace

import "github.com/infraboard/mcube/v2/ioc"

const (
	AppName = "trace"
)

func Get() *Trace {
	return ioc.Config().Get(AppName).(*Trace)
}
