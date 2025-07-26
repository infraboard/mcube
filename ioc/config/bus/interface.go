package bus

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
)

const (
	APP_NAME = "bus"
)

func GetService() Service {
	return ioc.Config().Get(APP_NAME).(Service)
}

type Service interface {
	Publisher
	SubScriber
}

type Publisher interface {
	// 事件发送
	Publish(ctx context.Context, e *Event) error
}

type SubScriber interface {
	// 主题订阅, 应用的多个实例 都会收到一份消息(广播)
	TopicSubscribe(ctx context.Context, subject string, cb EventHandler) error
	// 队列订阅, 默认应用名称为队列名称, 同一个队列中 只能收到一份消息
	QueueSubscribe(ctx context.Context, subject string, cb EventHandler) error
}

type EventHandler func(*Event)
