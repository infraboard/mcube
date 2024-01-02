package cors

import "github.com/infraboard/mcube/v2/ioc"

const (
	AppName = "cors"
)

func Get() *CORS {
	return ioc.Config().Get(AppName).(*CORS)
}
