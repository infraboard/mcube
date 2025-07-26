package kafka

import (
	"context"
	"os"
	"sync"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/bus"
	ioc_kafka "github.com/infraboard/mcube/v2/ioc/config/kafka"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/rs/zerolog"
	kafka "github.com/segmentio/kafka-go"
)

func init() {
	ioc.Config().Registry(&BusServiceImpl{
		producer: map[string]*kafka.Writer{},
		consumer: map[string]*kafka.Reader{},
		Group:    "mcube_bus",
	})
}

var _ bus.Service = (*BusServiceImpl)(nil)

type BusServiceImpl struct {
	ioc.ObjectImpl
	log *zerolog.Logger

	Group string `toml:"group" json:"group" yaml:"group"  env:"GROUP"`

	sync.Mutex
	hostname string
	producer map[string]*kafka.Writer
	consumer map[string]*kafka.Reader
}

func (b *BusServiceImpl) Name() string {
	return bus.APP_NAME
}

func (b *BusServiceImpl) Init() error {
	b.log = log.Sub(b.Name())
	b.hostname, _ = os.Hostname()
	return nil
}

func (b *BusServiceImpl) Close(ctx context.Context) {
	for _, p := range b.producer {
		p.Close()
	}

	for _, c := range b.consumer {
		c.Close()
	}
}

func (b *BusServiceImpl) GetProducer(topic string) *kafka.Writer {
	b.Lock()
	defer b.Unlock()

	if _, ok := b.producer[topic]; !ok {
		p := ioc_kafka.Get().Producer(topic)
		b.producer[topic] = p
	}
	return b.producer[topic]
}

func (b *BusServiceImpl) GetConsumer(group, topic string) *kafka.Reader {
	b.Lock()
	defer b.Unlock()

	if _, ok := b.consumer[topic]; !ok {
		p := ioc_kafka.Get().ConsumerGroup(group, []string{topic})
		b.consumer[topic] = p
	}
	return b.consumer[topic]
}

// 事件发送
func (b *BusServiceImpl) Publish(ctx context.Context, e *bus.Event) error {
	// Convert map[string][]string to []kafka.Header
	var headers []kafka.Header
	for k, vs := range e.Header {
		for _, v := range vs {
			headers = append(headers, kafka.Header{
				Key:   k,
				Value: []byte(v),
			})
		}
	}
	return b.GetProducer(e.Subject).WriteMessages(ctx, kafka.Message{
		Headers: headers,
		Value:   e.Data,
	})
}

// 订阅事件
func (b *BusServiceImpl) Subscribe(ctx context.Context, subject string, cb bus.EventHandler) error {
	for {
		m, err := b.GetConsumer(b.hostname, subject).ReadMessage(ctx)
		if err != nil {
			return err
		}

		// 打印日志
		b.log.Debug().Msgf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		// 事件转换
		// Convert []kafka.Header to map[string][]string
		headerMap := make(map[string][]string)
		for _, h := range m.Headers {
			headerMap[h.Key] = append(headerMap[h.Key], string(h.Value))
		}

		cb(&bus.Event{
			Subject: m.Topic,
			Header:  headerMap,
			Data:    m.Value,
		})
	}
}

// 订阅队列
func (b *BusServiceImpl) Queue(ctx context.Context, subject string, cb bus.EventHandler) error {
	for {
		m, err := b.GetConsumer(b.Group, subject).ReadMessage(ctx)
		if err != nil {
			return err
		}

		// 打印日志
		b.log.Debug().Msgf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		// 事件转换
		// Convert []kafka.Header to map[string][]string
		headerMap := make(map[string][]string)
		for _, h := range m.Headers {
			headerMap[h.Key] = append(headerMap[h.Key], string(h.Value))
		}

		cb(&bus.Event{
			Subject: m.Topic,
			Header:  headerMap,
			Data:    m.Value,
		})
	}
}
