package kafka

import (
	"errors"
	"sync"

	"github.com/infraboard/mcube/bus"
	"github.com/infraboard/mcube/bus/event"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"github.com/Shopify/sarama"
)

// NewBroker kafka broker
func NewBroker(conf *Config) (*Broker, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	kc, err := newSaramaConfig(conf)
	if err != nil {
		return nil, err
	}

	b := &Broker{
		conf: conf,
		kc:   kc,
		l:    zap.L().Named("Kafka Bus"),
	}

	return b, nil
}

// Broker kafka broker
type Broker struct {
	l    logger.Logger
	conf *Config
	kc   *sarama.Config

	producer sarama.AsyncProducer
	pubChan  chan<- *sarama.ProducerMessage
	mux      sync.Mutex
	wg       sync.WaitGroup
}

// Debug 日志
func (b *Broker) Debug(l logger.Logger) {
	b.l = l
}

// Connect 连接
func (b *Broker) Connect() error {
	b.mux.Lock()
	defer b.mux.Unlock()

	b.l.Debugf("connect: %v", b.conf.Hosts)

	// try to connect
	producer, err := sarama.NewAsyncProducer(b.conf.Hosts, b.kc)
	if err != nil {
		b.l.Errorf("Kafka connect fails with: %+v", err)
		return err
	}

	b.producer = producer
	b.pubChan = b.producer.Input()
	return nil
}

// Disconnect 端口连接
func (b *Broker) Disconnect() error {
	if b.producer != nil {
		if err := b.producer.Close(); err != nil {
			b.l.Errorf("Failed to close Kafka producer cleanly:", err)
		}
	}

	return nil
}

// Pub 发布事件
func (b *Broker) Pub(topic string, e *event.Event) error {
	if b.producer == nil || b.pubChan == nil {
		return errors.New("not connected")
	}

	return nil
}

// Sub 订阅事件
func (b *Broker) Sub(topic string, h bus.EventHandler) error {
	return nil
}
