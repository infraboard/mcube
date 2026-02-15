package ioc

import "context"

type Store interface {
	StoreUser
	StoreManage
}

type StoreUser interface {
	// 对象注册，返回自身支持链式调用
	Registry(obj Object) StoreUser
	// 批量注册对象
	RegistryAll(objs ...Object) StoreUser
	// 对象获取
	Get(name string, opts ...GetOption) Object
	// 根据对象类型, 直接加载对象
	Load(obj any, opts ...GetOption) error
	// 打印对象列表
	List() []string
	// 数量统计
	Len() int
	// 遍历注入的对象
	ForEach(fn func(*ObjectWrapper))
}

type StoreManage interface {
	// 从环境变量中加载对象配置
	LoadFromEnv(prefix string) error
}

// Object接口, 需要注册到ioc空间托管的对象需要实现的方法
type Object interface {
	// 对象初始化, 初始化对象的属性
	Init() error
	// 对象的名称, 根据名称可以从空间中取出对象
	Name() string
	// 对象版本, 默认v1
	Version() string
	// 对象优先级, 根据优先级 控制对象初始化的顺序
	Priority() int
	// 对象的销毁, 服务关闭时调用
	Close(ctx context.Context)
	// 对象一些元数据, 对象的更多描述信息, 扩展使用
	Meta() ObjectMeta
}

// 生命周期钩子接口 - 所有钩子都是可选的，按需实现

// PostConfigHook 配置加载完成后的钩子（在依赖注入之前）
// 适用场景：配置验证、配置预处理
type PostConfigHook interface {
	Object
	// OnPostConfig 配置加载完成后调用
	// 返回error会中断初始化流程
	OnPostConfig() error
}

// PreInitHook 对象初始化前的钩子
// 适用场景：准备工作、资源预分配
type PreInitHook interface {
	Object
	// OnPreInit Init()之前调用
	// 返回error会中断初始化流程
	OnPreInit() error
}

// PostInitHook 对象初始化后的钩子
// 适用场景：启动后台任务、注册监听器
type PostInitHook interface {
	Object
	// OnPostInit Init()成功后调用
	// 返回error会记录但不会中断流程
	OnPostInit() error
}

// PreStopHook 对象停止前的钩子
// 适用场景：优雅停机检查、等待请求完成
type PreStopHook interface {
	Object
	// OnPreStop Close()之前调用
	// context用于控制等待超时
	// 返回error会记录但不会中断关闭流程
	OnPreStop(ctx context.Context) error
}

// PostStopHook 对象停止后的钩子
// 适用场景：最终清理、资源释放确认
type PostStopHook interface {
	Object
	// OnPostStop Close()之后调用
	// context用于控制清理超时
	// 返回error会记录但不会中断流程
	OnPostStop(ctx context.Context) error
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
