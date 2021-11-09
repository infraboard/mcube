package app

import "fmt"

// InternalApp 内部服务实例, 不需要暴露
type InternalApp interface {
	Config() error
	Name() string
}

var (
	internalApps = map[string]InternalApp{}
)

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

// LoadGrpcApp 加载所有的Grpc app
func LoadInternalApp() error {
	for name, app := range internalApps {
		err := app.Config()
		if err != nil {
			return fmt.Errorf("config internal app %s error %s", name, err)
		}
	}
	return nil
}
