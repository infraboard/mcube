package pkg

// HTTPTemplate todo
const HTTPTemplate = `package pkg

import (
	"strings"

	"github.com/infraboard/mcube/http/router"
)

const (
	// ActionLableKey key name
	ActionLableKey = "action"
)

var (
	// GetAction Label
	GetAction = action("get")
	// ListAction label
	ListAction = action("list")
	// CreateAction label
	CreateAction = action("create")
	// UpdateAction label
	UpdateAction = action("update")
	// DeleteAction label
	DeleteAction = action("delete")
)
var (
	v1httpAPIs = make(map[string]HTTPAPI)
)

func action(value string) *router.Label {
	return router.NewLable(ActionLableKey, value)
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
