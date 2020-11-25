package kafka

import (
	"errors"
	"fmt"
	"sync"

	"github.com/infraboard/mcube/bus/event"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"github.com/Shopify/sarama"
)

// NewPublisher kafka broker
func NewPublisher(conf *PublisherConfig) (*Publisher, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	kc, err := newSaramaPubConfig(conf)
	if err != nil {
		return nil, err
	}

	b := &Publisher{
		conf: conf,
		kc:   kc,
		l:    zap.L().Named("Kafka Bus"),
	}

	return b, nil
}

// Publisher kafka broker
type Publisher struct {
	l    logger.Logger
	conf *PublisherConfig
	kc   *sarama.Config

	producer   sarama.AsyncProducer
	pubChan    chan<- *sarama.ProducerMessage
	pubSuccess chan<- *sarama.ProducerMessage
	pubFailed  chan<- *sarama.ProducerMessage
	mux        sync.Mutex
	wg         sync.WaitGroup
}

// Debug 日志
func (b *Publisher) Debug(l logger.Logger) {
	b.l = l
}

// Connect 连接
func (b *Publisher) Connect() error {
	b.mux.Lock()
	defer b.mux.Unlock()

	b.l.Debugf("connect: %v", b.conf.Hosts)

	// try to connect
	producer, err := sarama.NewAsyncProducer(b.conf.Hosts, b.kc)
	if err != nil {
		b.l.Errorf("new kafka producer fails with: %+v", err)
		return err
	}

	b.producer = producer
	b.pubChan = producer.Input()
	go b.watchSuccess(producer.Successes())
	go b.watchFailed(producer.Errors())

	return nil
}

// Disconnect 端口连接
func (b *Publisher) Disconnect() error {
	if b.producer != nil {
		if err := b.producer.Close(); err != nil {
			b.l.Errorf("Failed to close Kafka producer cleanly:", err)
		}
	}

	return nil
}

// Pub 发布事件
func (b *Publisher) Pub(topic string, e *event.Event) error {
	if e == nil {
		return fmt.Errorf("event is nil")
	}

	if err := e.Validate(); err != nil {
		return fmt.Errorf("validate event error, %s", err)
	}

	if b.producer == nil || b.pubChan == nil {
		return errors.New("not connected")
	}

	msg, err := newProducerMessage(e)
	if err != nil {
		return fmt.Errorf("new product message from event error, %s", err)
	}

	msg.Topic = topic
	b.pubChan <- msg
	return nil
}

func (b *Publisher) watchSuccess(ch <-chan *sarama.ProducerMessage) {
	for msg := range ch {
		b.l.Debugf("[%s] send mssage success, partition: %d, offset: %d", msg.Topic, msg.Partition, msg.Offset)
	}
}

func (b *Publisher) watchFailed(ch <-chan *sarama.ProducerError) {
	for msg := range ch {
		b.l.Errorf("[%s], send msg failed, %s", msg.Msg.Topic, msg.Err)
	}
}
