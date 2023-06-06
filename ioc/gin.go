package ioc

import (
	"fmt"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	ginApis = map[string]GinApiObject{}
)

// HTTPService Http服务的实例
type GinApiObject interface {
	IocObject
	Registry(gin.IRouter)
	Version() string
}

// RegistryGinApi 服务实例注册
func RegistryGinApi(obj GinApiObject) {
	// 已经注册的服务禁止再次注册
	_, ok := ginApis[obj.Name()]
	if ok {
		panic(fmt.Sprintf("gin app %s has registed", obj.Name()))
	}

	ginApis[obj.Name()] = obj
}

// LoadedGinApi 查询加载成功的服务
func LoadedGinApis() (apps []string) {
	for k := range ginApis {
		apps = append(apps, k)
	}
	return
}

func GetGinApi(name string) GinApiObject {
	app, ok := ginApis[name]
	if !ok {
		panic(fmt.Sprintf("http app %s not registed", name))
	}

	return app
}

// LoadGinApi 装载所有的gin app
func LoadGinApi(pathPrefix string, root gin.IRouter) {
	for _, api := range ginApis {
		if pathPrefix != "" && !strings.HasPrefix(pathPrefix, "/") {
			pathPrefix = "/" + pathPrefix
		}
		api.Registry(root.Group(path.Join(pathPrefix, api.Version(), api.Name())))
	}
}
