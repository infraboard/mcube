package gorestful

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/http/restful/accessor/form"
	"github.com/infraboard/mcube/v2/http/restful/accessor/yaml"
	"github.com/infraboard/mcube/v2/http/restful/accessor/yamlk8s"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/http"
)

const (
	AppName = "go_restful_webframework"
)

func Get() *restful.Container {
	return ioc.Config().Get(AppName).(*GoRestfulFramework).Container
}

func NewGoRestfulApiService(obj ioc.Object) *restful.WebService {
	ws := new(restful.WebService)
	ws.
		Path(http.Get().ApiObjectPathPrefix(obj)).
		Consumes(restful.MIME_JSON, form.MIME_POST_FORM, form.MIME_MULTIPART_FORM, yaml.MIME_YAML, yamlk8s.MIME_YAML).
		Produces(restful.MIME_JSON, yaml.MIME_YAML, yamlk8s.MIME_YAML)

	// 添加到Root Container
	Get().Add(ws)
	return ws
}
