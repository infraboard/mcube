package ioc

// IocObject 内部服务实例, 不需要暴露
type IocObject interface {
	// 对象初始化
	Init() error
	// 对象的名称
	Name() string
	// 对象版本
	Version() string
	// 对象优先级
	Priority() int
}

const (
	DEFAULT_VERSION = "v1"
)

type IocObjectImpl struct {
}

func (i *IocObjectImpl) Init() error {
	return nil
}

func (i *IocObjectImpl) Name() string {
	return ""
}

func (i *IocObjectImpl) Version() string {
	return DEFAULT_VERSION
}

func (i *IocObjectImpl) Priority() int {
	return 0
}
