package rabbitmq

const (
	EXCHANGE_TYPE_DIRECT  EXCHANGE_TYPE = "direct"  // Routing Key 完全匹配 Binding Key - 精确路由（如任务分发）
	EXCHANGE_TYPE_FANOUT  EXCHANGE_TYPE = "fanout"  // 忽略 Routing Key，广播到所有队列 - 发布/订阅（如事件通知）
	EXCHANGE_TYPE_TOPIC   EXCHANGE_TYPE = "topic"   // Routing Key 通配符匹配 Binding Key - 动态路由（如分类日志）
	EXCHANGE_TYPE_HEADERS EXCHANGE_TYPE = "headers" // 匹配消息的 Headers 属性 - 复杂条件过滤（如协议适配）
	EXCHANGE_TYPE_DEFAULT EXCHANGE_TYPE = ""        // 默认交换机: Routing Key = 队列名 - 简单测试或直接队列通信
)

type EXCHANGE_TYPE string

func (e EXCHANGE_TYPE) String() string {
	return string(e)
}
