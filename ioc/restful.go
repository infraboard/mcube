package ioc

import (
	"fmt"
	"strings"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/http/restful/accessor/form"
	"github.com/infraboard/mcube/http/restful/accessor/yaml"
	"github.com/infraboard/mcube/http/restful/accessor/yamlk8s"
)

var (
	goRestfulApis = map[string]GoRestfulApiObject{}
)

// HTTPService Http服务的实例
type GoRestfulApiObject interface {
	IocObject
	Registry(*restful.WebService)
	Version() string
}

// RegistryGoRestfulApi 服务实例注册
func RegistryGoRestfulApi(obj GoRestfulApiObject) {
	// 已经注册的服务禁止再次注册
	_, ok := goRestfulApis[obj.Name()]
	if ok {
		panic(fmt.Sprintf("http app %s has registed", obj.Name()))
	}

	goRestfulApis[obj.Name()] = obj
}

// LoadedHttpApp 查询加载成功的服务
func LoadedGoRestfulApi() (names []string) {
	for k := range goRestfulApis {
		names = append(names, k)
	}
	return
}

func GetGoRestfulApi(name string) GoRestfulApiObject {
	app, ok := goRestfulApis[name]
	if !ok {
		panic(fmt.Sprintf("http app %s not registed", name))
	}

	return app
}

// LoadHttpApp 装载所有的http app
func LoadGoRestfulApi(pathPrefix string, root *restful.Container) {
	for _, api := range goRestfulApis {
		pathPrefix = strings.TrimSuffix(pathPrefix, "/")
		ws := new(restful.WebService)
		ws.
			Path(fmt.Sprintf("%s/%s/%s", pathPrefix, api.Version(), api.Name())).
			Consumes(restful.MIME_JSON, form.MIME_POST_FORM, form.MIME_MULTIPART_FORM, yaml.MIME_YAML, yamlk8s.MIME_YAML).
			Produces(restful.MIME_JSON, yaml.MIME_YAML, yamlk8s.MIME_YAML)
		api.Registry(ws)
		root.Add(ws)
	}
}
