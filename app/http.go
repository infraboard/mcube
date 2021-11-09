package app

import (
	"fmt"
	"strings"

	"github.com/infraboard/mcube/http/router"
)

var (
	httpApps = map[string]HTTPApp{}
)

// HTTPService Http服务的实例
type HTTPApp interface {
	Registry(router.SubRouter)
	Config() error
	Name() string
}

// RegistryHttpApp 服务实例注册
func RegistryHttpApp(app HTTPApp) {
	// 已经注册的服务禁止再次注册
	_, ok := httpApps[app.Name()]
	if ok {
		panic(fmt.Sprintf("http app %s has registed", app.Name()))
	}

	httpApps[app.Name()] = app
}

// LoadedHttpApp 查询加载成功的服务
func LoadedHttpApp() (apps []string) {
	for k := range httpApps {
		apps = append(apps, k)
	}
	return
}

func GetHttpApp(name string) HTTPApp {
	app, ok := httpApps[name]
	if !ok {
		panic(fmt.Sprintf("http app %s not registed", name))
	}

	return app
}

// LoadHttpApp 装载所有的http app
func LoadHttpApp(pathPrefix string, root router.Router) error {
	for _, api := range httpApps {
		if err := api.Config(); err != nil {
			return err
		}
		if pathPrefix != "" && !strings.HasPrefix(pathPrefix, "/") {
			pathPrefix = "/" + pathPrefix
		}
		api.Registry(root.SubRouter(pathPrefix))
	}
	return nil
}
