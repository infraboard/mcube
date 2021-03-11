package pkg

// ServiceTemplate todo
const ServiceTemplate = `package pkg

import (
	"fmt"

	"github.com/infraboard/mcube/pb/http"
	"google.golang.org/grpc"
)

var (
	servers       []Service
	successLoaded []string

	entrySet = http.NewEntrySet()
)

// InitV1GRPCAPI 初始化GRPC服务
func InitV1GRPCAPI(server *grpc.Server) {
	return
}

// HTTPEntry todo
func HTTPEntry() *http.EntrySet {
	return entrySet
}

// GetGrpcPathEntry todo
func GetGrpcPathEntry(path string) *http.Entry {
	es := HTTPEntry()
	for i := range es.Items {
		if es.Items[i].GrpcPath == path {
			return es.Items[i]
		}
	}

	return nil
}

// LoadedService 查询加载成功的服务
func LoadedService() []string {
	return successLoaded
}
func addService(name string, svr Service) {
	servers = append(servers, svr)
	successLoaded = append(successLoaded, name)
}

// Service 注册上的服务必须实现的方法
type Service interface {
	Config() error
	HTTPEntry() *http.EntrySet
}

// RegistryService 服务实例注册
func RegistryService(name string, svr Service) {
	switch value := svr.(type) {
	default:
		fmt.Println(value)
		panic(fmt.Sprintf("unknown service type %s", name))
	}
}

func registryError(name string) {
	panic("service " + name + " has registried")
}

// InitService 初始化所有服务
func InitService() error {
	for _, s := range servers {
		if err := s.Config(); err != nil {
			return err
		}
		entrySet.Merge(s.HTTPEntry())
	}
	return nil
}`
