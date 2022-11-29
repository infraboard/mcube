package app

import (
	"fmt"
	"strings"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/http/restful/accessor/form"
	"github.com/infraboard/mcube/http/restful/accessor/yaml"
	"github.com/infraboard/mcube/http/restful/accessor/yamlk8s"
)

var (
	restfulApps = map[string]RESTfulApp{}
)

// HTTPService Http服务的实例
type RESTfulApp interface {
	InternalApp
	Registry(*restful.WebService)
	Version() string
}

// RegistryRESTfulApp 服务实例注册
func RegistryRESTfulApp(app RESTfulApp) {
	// 已经注册的服务禁止再次注册
	_, ok := restfulApps[app.Name()]
	if ok {
		panic(fmt.Sprintf("http app %s has registed", app.Name()))
	}

	restfulApps[app.Name()] = app
}

// LoadedHttpApp 查询加载成功的服务
func LoadedRESTfulApp() (apps []string) {
	for k := range restfulApps {
		apps = append(apps, k)
	}
	return
}

func GetRESTfulApp(name string) RESTfulApp {
	app, ok := restfulApps[name]
	if !ok {
		panic(fmt.Sprintf("http app %s not registed", name))
	}

	return app
}

// LoadHttpApp 装载所有的http app
func LoadRESTfulApp(pathPrefix string, root *restful.Container) {
	for _, api := range restfulApps {
		pathPrefix = strings.TrimSuffix(pathPrefix, "/")
		ws := new(restful.WebService)
		ws.
			Path(fmt.Sprintf("%s/%s/%s", pathPrefix, api.Version(), api.Name())).
			Consumes(restful.MIME_JSON, restful.MIME_XML, form.MIME_POST_FORM, form.MIME_MULTIPART_FORM, yaml.MIME_YAML, yamlk8s.MIME_YAML).
			Produces(restful.MIME_JSON, restful.MIME_XML, yaml.MIME_YAML, yamlk8s.MIME_YAML)
		api.Registry(ws)
		root.Add(ws)
	}
}
