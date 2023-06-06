package ioc

import (
	"fmt"
	"strings"

	"github.com/infraboard/mcube/http/router"
)

var (
	httpApis = map[string]HTTPApiObject{}
)

// HTTPService Http服务的实例
type HTTPApiObject interface {
	IocObject
	Registry(router.SubRouter)
}

// RegistryHttpApi 服务实例注册
func RegistryHttpApi(obj HTTPApiObject) {
	// 已经注册的服务禁止再次注册
	_, ok := httpApis[obj.Name()]
	if ok {
		panic(fmt.Sprintf("http app %s has registed", obj.Name()))
	}

	httpApis[obj.Name()] = obj
}

// LoadedHttpApi 查询加载成功的服务
func LoadedHttpApis() (apps []string) {
	for k := range httpApis {
		apps = append(apps, k)
	}
	return
}

func GetHttpApi(name string) HTTPApiObject {
	app, ok := httpApis[name]
	if !ok {
		panic(fmt.Sprintf("http app %s not registed", name))
	}

	return app
}

// LoadHttpApi 装载所有的http app
func LoadHttpApi(pathPrefix string, root router.Router) {
	for _, api := range httpApis {
		if pathPrefix != "" && !strings.HasPrefix(pathPrefix, "/") {
			pathPrefix = "/" + pathPrefix
		}
		api.Registry(root.SubRouter(pathPrefix))
	}
}
