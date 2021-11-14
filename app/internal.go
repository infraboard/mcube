package app

import "fmt"

var (
	internalApps = map[string]InternalApp{}
)

// InternalApp 内部服务实例, 不需要暴露
type InternalApp interface {
	Config() error
	Name() string
}

// RegistryInternalApp 服务实例注册
func RegistryInternalApp(app InternalApp) {
	// 已经注册的服务禁止再次注册
	_, ok := internalApps[app.Name()]
	if ok {
		panic(fmt.Sprintf("internal app %s has registed", app.Name()))
	}

	internalApps[app.Name()] = app
}

// LoadedInternalApp 查询加载成功的服务
func LoadedInternalApp() (apps []string) {
	for k := range internalApps {
		apps = append(apps, k)
	}
	return
}

func GetInternalApp(name string) InternalApp {
	app, ok := internalApps[name]
	if !ok {
		panic(fmt.Sprintf("internal app %s not registed", name))
	}

	return app
}
