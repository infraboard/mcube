package rabbitmq

import (
	"context"
	"fmt"
	"sync"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	"github.com/infraboard/mcube/v2/ioc/config/bus"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/mcube/v2/ioc/config/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Config().Registry(&BusServiceImpl{
		publishers: map[string]*rabbitmq.Publisher{},
		consumers:  map[string]*rabbitmq.Consumer{},
	})
}

var _ bus.Service = (*BusServiceImpl)(nil)

type BusServiceImpl struct {
	ioc.ObjectImpl
	log *zerolog.Logger

	Group string `toml:"group" json:"group" yaml:"group"  env:"GROUP"`

	publishers map[string]*rabbitmq.Publisher
	consumers  map[string]*rabbitmq.Consumer

	mu sync.Mutex
}

func (b *BusServiceImpl) Name() string {
	return bus.APP_NAME
}

func (b *BusServiceImpl) Init() error {
	if b.Group == "" {
		b.Group = application.Get().GetAppName()
	}

	b.log = log.Sub(b.Name())
	return nil
}

func (b *BusServiceImpl) Close(ctx context.Context) {
	for k, v := range b.publishers {
		v.Close()
		b.log.Info().Msgf("close %s publisher", k)
	}

	for k, v := range b.consumers {
		v.Close()
		b.log.Info().Msgf("close %s publisher", k)
	}

}

func (b *BusServiceImpl) GetPublisher(subject string) (*rabbitmq.Publisher, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if v, ok := b.publishers[subject]; ok {
		return v, nil
	}

	publisher, err := rabbitmq.NewPublisher()
	if err != nil {
		return nil, err
	}
	b.publishers[subject] = publisher

	return publisher, nil
}

func (b *BusServiceImpl) GetConsumer(subject string) (*rabbitmq.Consumer, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if v, ok := b.consumers[subject]; ok {
		return v, nil
	}

	consumer, err := rabbitmq.NewConsumer()
	if err != nil {
		return nil, err
	}
	b.consumers[subject] = consumer

	return consumer, nil
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

	p, err := b.GetPublisher(e.Subject)
	if err != nil {
		return err
	}

	return p.Publish(ctx, msg)
}

// 订阅逻辑（广播模式）
func (b *BusServiceImpl) TopicSubscribe(ctx context.Context, subject string, cb bus.EventHandler) error {
	consumer, err := b.GetConsumer(subject)
	if err != nil {
		return err
	}

	// 使用随机队列名 + 绑定到 Topic Exchange
	err = consumer.TopicSubscribe(ctx, b.Group, subject, func(ctx context.Context, msg *rabbitmq.Message) error {
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
func (b *BusServiceImpl) QueueSubscribe(ctx context.Context, queue string, cb bus.EventHandler) error {
	consumer, err := b.GetConsumer(queue)
	if err != nil {
		return err
	}
	// 使用固定队列名 + 绑定到 Topic Exchange
	err = consumer.TopicSubscribe(ctx, b.Group, queue, func(ctx context.Context, msg *rabbitmq.Message) error {
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
