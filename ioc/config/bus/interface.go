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
	// 订阅事件
	Subscribe(ctx context.Context, subject string, cb EventHandler) error
	// 订阅队列
	// subject: 主题
	// queue: 队列名称, kafka则为消费组名称
	Queue(ctx context.Context, subject string, cb EventHandler) error
}

type EventHandler func(*Event)
