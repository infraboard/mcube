package pkg

// HTTPTemplate todo
const HTTPTemplate = `package pkg

import (
	"strings"

	"github.com/infraboard/mcube/http/router"
)

var (
	v1httpAPIs = make(map[string]HTTPAPI)
)

// LoadedHTTP 查询加载成功的HTTP API
func LoadedHTTP() []string {
	var apis []string
	for k := range v1httpAPIs {
		apis = append(apis, k)
	}
	return apis
}

// HTTPAPI restful 服务
type HTTPAPI interface {
	Registry(router router.SubRouter)
	Config() error
}

// RegistryHTTPV1 注册HTTP服务
func RegistryHTTPV1(name string, api HTTPAPI) {
	if _, ok := v1httpAPIs[name]; ok {
		panic("http api " + name + " has registry")
	}
	v1httpAPIs[name] = api
}

// InitV1HTTPAPI 初始化API服务
func InitV1HTTPAPI(pathPrefix string, root router.Router) error {
	for _, api := range v1httpAPIs {
		if err := api.Config(); err != nil {
			return err
		}
		if pathPrefix != "" && !strings.HasPrefix(pathPrefix, "/") {
			pathPrefix = "/" + pathPrefix
		}
		api.Registry(root.SubRouter(pathPrefix + "/api/v1"))
	}
	return nil
}
`
