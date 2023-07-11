package ioc

import (
	"fmt"
	"path"
	"strings"

	"github.com/emicklei/go-restful/v3"
	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/http/restful/accessor/form"
	"github.com/infraboard/mcube/http/restful/accessor/yaml"
	"github.com/infraboard/mcube/http/restful/accessor/yamlk8s"
)

const (
	ApiNamespace = "apis"
)

type GinApiObject interface {
	IocObject
	Registry(gin.IRouter)
}

type GoRestfulApiObject interface {
	IocObject
	Registry(*restful.WebService)
}

// 注册API对象
func RegistryApi(obj IocObject) {
	RegistryObjectWithNs(ApiNamespace, obj)
}

// 获取API对象
func GetApi(name string) IocObject {
	return GetObjectWithNs(ApiNamespace, name, DEFAULT_VERSION)
}

// 查询已经注册的API对象名称
func ListApiObjectNames() (names []string) {
	return store.Namespace(ApiNamespace).ObjectUids()
}

// LoadGinApi 装载所有的gin app
func LoadGinApi(pathPrefix string, root gin.IRouter) {
	objects := store.Namespace(ApiNamespace)
	for _, obj := range objects.Items {
		api, ok := obj.(GinApiObject)
		if !ok {
			continue
		}

		if pathPrefix != "" && !strings.HasPrefix(pathPrefix, "/") {
			pathPrefix = "/" + pathPrefix
		}
		api.Registry(root.Group(path.Join(pathPrefix, api.Version(), api.Name())))
	}
}

// LoadHttpApp 装载所有的http app
func LoadGoRestfulApi(pathPrefix string, root *restful.Container) {
	objects := store.Namespace(ApiNamespace)
	for _, obj := range objects.Items {
		api, ok := obj.(GoRestfulApiObject)
		if !ok {
			continue
		}

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
