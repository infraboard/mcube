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
		Exchange: "mcube_bus",
	})
}

var _ bus.Service = (*BusServiceImpl)(nil)

type BusServiceImpl struct {
	ioc.ObjectImpl

	Exchange string `toml:"exchange" json:"exchange" yaml:"exchange"  env:"EXCHANGE"`

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

// 事件发送
func (b *BusServiceImpl) Publish(ctx context.Context, e *bus.Event) error {
	msg := &rabbitmq.Message{
		Exchange:   b.Exchange,
		RoutingKey: e.Subject,
		Body:       e.Data,
		Headers:    make(amqp091.Table),
	}

	for k, v := range e.Header {
		msg.Headers[k] = v
	}

	return b.publisher.Publish(ctx, msg)
}

// 订阅事件
func (b *BusServiceImpl) Subscribe(ctx context.Context, subject string, cb bus.EventHandler) error {
	err := b.consumer.TopicSubscribe(ctx, b.Exchange, subject, func(ctx context.Context, msg *rabbitmq.Message) error {
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

// 订阅队列
func (b *BusServiceImpl) Queue(ctx context.Context, subject string, queue string, cb bus.EventHandler) error {
	err := b.consumer.DirectSubscribe(ctx, b.Exchange, queue, func(ctx context.Context, msg *rabbitmq.Message) error {
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
