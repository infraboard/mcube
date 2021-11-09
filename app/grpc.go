package app

import (
	"fmt"

	"google.golang.org/grpc"
)

// GRPCService GRPC服务的实例
type GRPCApp interface {
	Registry(*grpc.Server)
	Config() error
	Name() string
}

var (
	grpcApps = map[string]GRPCApp{}
)

// RegistryService 服务实例注册
func RegistryGrpcApp(app GRPCApp) {
	// 已经注册的服务禁止再次注册
	_, ok := grpcApps[app.Name()]
	if ok {
		panic(fmt.Sprintf("grpc app %s has registed", app.Name()))
	}

	grpcApps[app.Name()] = app
}

// LoadedGrpcApp 查询加载成功的服务
func LoadedGrpcApp() (apps []string) {
	for k := range grpcApps {
		apps = append(apps, k)
	}
	return
}

func GetGrpcApp(name string) GRPCApp {
	app, ok := grpcApps[name]
	if !ok {
		panic(fmt.Sprintf("grpc app %s not registed", name))
	}

	return app
}

// LoadGrpcApp 加载所有的Grpc app
func LoadGrpcApp(server *grpc.Server) error {
	for name, app := range grpcApps {
		err := app.Config()
		if err != nil {
			return fmt.Errorf("config grpc app %s error %s", name, err)
		}

		app.Registry(server)
	}
	return nil
}
