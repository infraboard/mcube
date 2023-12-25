package http

import "github.com/infraboard/mcube/v2/ioc"

const (
	AppName = "http"
)

func Get() *Http {
	return ioc.Config().Get(AppName).(*Http)
}
