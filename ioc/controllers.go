package ioc

import "fmt"

var (
	controllers = map[string]IocObject{}
)

// IocObject 内部服务实例, 不需要暴露
type IocObject interface {
	// 对象初始化
	Init() error
	// 对象的名称
	Name() string
}

// RegistryController 服务实例注册
func RegistryController(obj IocObject) {
	// 已经注册的服务禁止再次注册
	_, ok := controllers[obj.Name()]
	if ok {
		panic(fmt.Sprintf("internal app %s has registed", obj.Name()))
	}

	controllers[obj.Name()] = obj
}

// LoadedControllers 查询加载成功的服务
func LoadedControllers() (names []string) {
	for k := range controllers {
		names = append(names, k)
	}
	return
}

func GetController(name string) IocObject {
	obj, ok := controllers[name]
	if !ok {
		panic(fmt.Sprintf("internal obj %s not registed", name))
	}

	return obj
}
