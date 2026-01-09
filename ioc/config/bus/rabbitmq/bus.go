package rabbitmq

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
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

	// group 队列模式下的 队列名称或者消费组名称，一个组里面的实例消费一个队列
	Group string `toml:"group" json:"group" yaml:"group"  env:"GROUP"`
	// nodename 广播模式下的节点名称，默认hostname, 每个节点独立一个 节点队列: group.nodename
	NodeName string `toml:"node_name" json:"node_name" yaml:"node_name" env:"NODE_NAME"`

	publishers map[string]*rabbitmq.Publisher
	consumers  map[string]*rabbitmq.Consumer

	mu sync.Mutex
}

func (b *BusServiceImpl) Name() string {
	return bus.APP_NAME
}

func (i *BusServiceImpl) Priority() int {
	return bus.APP_PRIORITY
}

func (b *BusServiceImpl) Init() error {
	if b.Group == "" {
		b.Group = application.Get().GetAppName()
	}

	if b.NodeName == "" {
		hostname, err := os.Hostname()
		if err != nil {
			return err
		}

		b.NodeName = hostname
	}

	b.log = log.Sub(b.Name())
	return nil
}

func (b *BusServiceImpl) Close(ctx context.Context) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for k, v := range b.publishers {
		v.Close()
		b.log.Info().Msgf("close %s publisher", k)
	}

	for k, v := range b.consumers {
		v.Close()
		b.log.Info().Msgf("close %s consumer", k)
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

	// 使用group + nodename + 绑定到 Topic Exchange
	err = consumer.TopicSubscribe(ctx, b.Group, subject, func(ctx context.Context, msg *rabbitmq.Message) error {
		cb(&bus.Event{
			Subject: msg.RoutingKey,
			Header:  b.convert(msg.Headers),
			Data:    msg.Body,
		})
		return nil
	}, rabbitmq.WithDurable(true),
		rabbitmq.WithAutoDelete(true),
		rabbitmq.WithQueueName(sanitizeQueueName(b.Group+"."+subject+"."+b.NodeName)))
	if err != nil {
		return err
	}
	return nil
}

// 队列逻辑（竞争消费模式）
func (b *BusServiceImpl) QueueSubscribe(ctx context.Context, subject string, cb bus.EventHandler) error {
	consumer, err := b.GetConsumer(subject)
	if err != nil {
		return err
	}
	// 使用固定队列名 group + 绑定到 Topic Exchange
	err = consumer.TopicSubscribe(ctx, b.Group, subject, func(ctx context.Context, msg *rabbitmq.Message) error {
		cb(&bus.Event{
			Subject: msg.RoutingKey,
			Header:  b.convert(msg.Headers),
			Data:    msg.Body,
		})
		return nil
	}, rabbitmq.WithDurable(true),
		rabbitmq.WithAutoDelete(false),
		rabbitmq.WithQueueName(sanitizeQueueName(b.Group+"."+subject)))
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

// sanitizeQueueName 清理队列名称，使其符合 RabbitMQ 命名规范
// 只允许字母、数字、连字符、点号、下划线
func sanitizeQueueName(name string) string {
	// 替换非法字符为下划线
	reg := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	sanitized := reg.ReplaceAllString(name, "_")

	// 移除开头和结尾的点号
	sanitized = strings.Trim(sanitized, ".")

	// 确保不为空
	if sanitized == "" {
		sanitized = "default"
	}

	// 限制长度（RabbitMQ 队列名最多 255 字节）
	if len(sanitized) > 255 {
		sanitized = sanitized[:255]
		// 再次移除结尾的点号
		sanitized = strings.TrimRight(sanitized, ".")
	}

	return sanitized
}
