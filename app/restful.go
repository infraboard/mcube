package app

import (
	"fmt"
	"strings"

	"github.com/emicklei/go-restful/v3"
)

var (
	restfulApps = map[string]RESTfulApp{}
)

// HTTPService Http服务的实例
type RESTfulApp interface {
	Registry(*restful.WebService)
	Config() error
	Name() string
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
	for name, api := range restfulApps {
		if !strings.HasSuffix(pathPrefix, "/") {
			pathPrefix = pathPrefix + "/"
		}
		ws := new(restful.WebService)
		ws.
			Path(pathPrefix+name).
			Consumes(restful.MIME_XML, restful.MIME_JSON).
			Produces(restful.MIME_JSON, restful.MIME_XML)

		api.Registry(ws)
		root.Add(ws)
	}
}
