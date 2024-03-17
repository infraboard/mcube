package http

import "github.com/infraboard/mcube/v2/ioc"

const (
	AppName = "http"
)

func Get() *Http {
	obj := ioc.Config().Get(AppName)
	if obj == nil {
		return defaultConfig
	}
	return obj.(*Http)
}
