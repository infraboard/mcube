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
	GorestfulApiNamespace = "go-restful-v3"
)

// HTTPService Http服务的实例
type GoRestfulApiObject interface {
	IocObject
	// 注册的路由
	Registry(*restful.WebService)
	// Api版本
	Version() string
}

// RegistryGoRestfulApi 服务实例注册
func RegistryGoRestfulApi(obj GoRestfulApiObject) {
	RegistryObjectWithNs(GorestfulApiNamespace, obj)
}

// 查询注入的对象
func GetGoRestfulApi(name string) IocObject {
	return GetObjectWithNs(GorestfulApiNamespace, name)
}

// LoadedHttpApp 查询加载成功的服务
func LoadedGoRestfulApi() (names []string) {
	return store.Namespace(GorestfulApiNamespace).ObjectNames()
}

// LoadHttpApp 装载所有的http app
func LoadGoRestfulApi(pathPrefix string, root *restful.Container) {
	objects := store.Namespace(GorestfulApiNamespace)
	for _, obj := range objects.Items {
		api := obj.(GoRestfulApiObject)
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
