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
	API_NAMESPACE = "apis"
)

// 用于托管RestApi对象的Ioc空间, 最后初始化
func Api() StoreUser {
	return store.Namespace(API_NAMESPACE)
}

type GinApiObject interface {
	Object
	Registry(gin.IRouter)
}

type GoRestfulApiObject interface {
	Object
	Registry(*restful.WebService)
}

// 注册API对象
func RegistryApi(obj Object) {
	RegistryObjectWithNs(API_NAMESPACE, obj)
}

// 获取API对象
func GetApi(name string) Object {
	return GetObjectWithNs(API_NAMESPACE, name)
}

// LoadGinApi 装载所有的gin app
func LoadGinApi(pathPrefix string, root gin.IRouter) {
	objects := store.Namespace(API_NAMESPACE)
	objects.ForEach(func(w *ObjectWrapper) {
		api, ok := w.Value.(GinApiObject)
		if !ok {
			return
		}

		if pathPrefix != "" && !strings.HasPrefix(pathPrefix, "/") {
			pathPrefix = "/" + pathPrefix
		}
		api.Registry(root.Group(path.Join(pathPrefix, api.Version(), api.Name())))
	})
}

// LoadHttpApp 装载所有的http app
func LoadGoRestfulApi(pathPrefix string, root *restful.Container) {
	objects := store.Namespace(API_NAMESPACE)
	objects.ForEach(func(w *ObjectWrapper) {
		api, ok := w.Value.(GoRestfulApiObject)
		if !ok {
			return
		}

		pathPrefix = strings.TrimSuffix(pathPrefix, "/")
		ws := new(restful.WebService)
		ws.
			Path(fmt.Sprintf("%s/%s/%s", pathPrefix, api.Version(), api.Name())).
			Consumes(restful.MIME_JSON, form.MIME_POST_FORM, form.MIME_MULTIPART_FORM, yaml.MIME_YAML, yamlk8s.MIME_YAML).
			Produces(restful.MIME_JSON, yaml.MIME_YAML, yamlk8s.MIME_YAML)
		api.Registry(ws)
		root.Add(ws)
	})
}
