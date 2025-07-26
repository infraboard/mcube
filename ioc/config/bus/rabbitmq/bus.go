package rabbitmq

import (
	"context"
	"fmt"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/bus"
	"github.com/infraboard/mcube/v2/ioc/config/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
)

func init() {
	ioc.Config().Registry(&BusServiceImpl{
		Group: "mcube_bus",
	})
}

var _ bus.Service = (*BusServiceImpl)(nil)

type BusServiceImpl struct {
	ioc.ObjectImpl

	Group string `toml:"group" json:"group" yaml:"group"  env:"GROUP"`

	publisher *rabbitmq.Publisher
	consumer  *rabbitmq.Consumer
}

func (b *BusServiceImpl) Name() string {
	return bus.APP_NAME
}

func (b *BusServiceImpl) Init() error {
	publisher, err := rabbitmq.NewPublisher()
	if err != nil {
		return err
	}
	b.publisher = publisher

	consumer, err := rabbitmq.NewConsumer()
	if err != nil {
		return err
	}
	b.consumer = consumer
	return nil
}

func (i *BusServiceImpl) Close(ctx context.Context) {
	if i.publisher != nil {
		i.publisher.Close()
	}
	if i.consumer != nil {
		i.consumer.Close()
	}
}

// 发布逻辑（始终发布到 Topic Exchange）
func (b *BusServiceImpl) Publish(ctx context.Context, e *bus.Event) error {
	msg := &rabbitmq.Message{
		Exchange:   b.Group,   // 固定为 Topic Exchange
		RoutingKey: e.Subject, // 路由键 = 事件主题
		Body:       e.Data,
		Headers:    make(amqp091.Table),
	}

	for k, v := range e.Header {
		msg.Headers[k] = v
	}

	return b.publisher.Publish(ctx, msg)
}

// 订阅逻辑（广播模式）
func (b *BusServiceImpl) Subscribe(ctx context.Context, subject string, cb bus.EventHandler) error {
	// 使用随机队列名 + 绑定到 Topic Exchange
	err := b.consumer.TopicSubscribe(ctx, b.Group, subject, func(ctx context.Context, msg *rabbitmq.Message) error {
		cb(&bus.Event{
			Subject: msg.Exchange,
			Header:  b.convert(msg.Headers),
			Data:    msg.Body,
		})
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// 队列逻辑（竞争消费模式）
func (b *BusServiceImpl) Queue(ctx context.Context, queue string, cb bus.EventHandler) error {
	// 使用固定队列名 + 绑定到 Topic Exchange
	err := b.consumer.TopicSubscribe(ctx, b.Group, queue, func(ctx context.Context, msg *rabbitmq.Message) error {
		cb(&bus.Event{
			Subject: msg.RoutingKey,
			Header:  b.convert(msg.Headers),
			Data:    msg.Body,
		})
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *BusServiceImpl) convert(table amqp091.Table) map[string][]string {
	headers := make(map[string][]string)
	for k, v := range table {
		if str, ok := v.(string); ok {
			headers[k] = []string{str}
		} else {
			headers[k] = []string{fmt.Sprintf("%v", v)}
		}
	}
	return headers
}
