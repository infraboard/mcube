package ioc

import (
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	GinApiNamespace = "gin"
)

// HTTPService Http服务的实例
type GinApiObject interface {
	IocObject
	Registry(gin.IRouter)
	Version() string
}

// RegistryGinApi 服务实例注册
func RegistryGinApi(obj GinApiObject) {
	RegistryObjectWithNs(GinApiNamespace, obj)
}

// Gin实例名称
func GetGinApi(name string) IocObject {
	return GetObjectWithNs(GinApiNamespace, name)
}

// LoadedGinApi 查询加载成功的服务
func LoadedGinApis() (names []string) {
	return store.Namespace(GinApiNamespace).ObjectNames()
}

// LoadGinApi 装载所有的gin app
func LoadGinApi(pathPrefix string, root gin.IRouter) {
	objects := store.Namespace(GinApiNamespace)
	for _, obj := range objects.Items {
		api := obj.(GinApiObject)
		if pathPrefix != "" && !strings.HasPrefix(pathPrefix, "/") {
			pathPrefix = "/" + pathPrefix
		}
		api.Registry(root.Group(path.Join(pathPrefix, api.Version(), api.Name())))
	}
}
