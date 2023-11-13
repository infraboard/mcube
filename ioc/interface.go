package ioc

import "context"

type Stroe interface {
	StoreUser
	StoreManage
}

type StoreUser interface {
	// 对象注册
	Registry(obj Object)
	// 对象获取
	Get(name string, opts ...GetOption) Object
	// 根据对象类型, 直接加载对象
	Load(obj any, opts ...GetOption) error
	// 打印对象列表
	List() []string
	// 数量统计
	Count() int
	// 遍历注入的对象
	ForEach(fn func(*ObjectWrapper))
}

type StoreManage interface {
	// 从环境变量中加载对象配置
	LoadFromEnv(prefix string) error
}

// Object 内部服务实例, 不需要暴露
type Object interface {
	// 对象初始化
	Init() error
	// 对象的名称
	Name() string
	// 对象版本
	Version() string
	// 对象优先级
	Priority() int
	// 对象的销毁
	Close(ctx context.Context) error
	// 是否允许同名对象被替换, 默认不允许被替换
	AllowOverwrite() bool
}

func WithVersion(v string) GetOption {
	return func(o *option) {
		o.version = v
	}
}

type GetOption func(*option)

func defaultOption() *option {
	return &option{
		version: DEFAULT_VERSION,
	}
}

type option struct {
	version string
}

func (o *option) Apply(opts ...GetOption) *option {
	for i := range opts {
		opts[i](o)
	}
	return o
}
