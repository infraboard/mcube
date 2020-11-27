package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/Shopify/sarama"

	"github.com/infraboard/mcube/bus"
	"github.com/infraboard/mcube/bus/event"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
)

// NewSubscriber kafka broker
func NewSubscriber(conf *Config) (*Subscriber, error) {
	if err := conf.ValidateSubscriberConfig(); err != nil {
		return nil, err
	}

	kc, err := newSaramaSubConfig(conf.baseConfig, conf.subscriberConfig)
	if err != nil {
		return nil, err
	}

	b := &Subscriber{
		conf:  conf,
		kc:    kc,
		l:     zap.L().Named("Kafka Bus"),
		ready: make(chan bool),
	}

	return b, nil
}

// Subscriber kafka broker
type Subscriber struct {
	l    logger.Logger
	conf *Config
	kc   *sarama.Config
	h    bus.EventHandler

	comsummer sarama.ConsumerGroup
	mux       sync.Mutex
	ready     chan bool
	cancel    context.CancelFunc
}

// Connect 连接
func (s *Subscriber) Connect() error {
	s.mux.Lock()
	defer s.mux.Unlock()

	// try to connect
	s.l.Debugf("try connect: %v ...", s.conf.Hosts)
	consumer, err := sarama.NewConsumerGroup(s.conf.Hosts, s.conf.GroupID, s.kc)
	if err != nil {
		s.l.Errorf("Kafka consummer connect fails with: %+v", err)
		return err
	}
	s.comsummer = consumer
	s.l.Debugf("connect %v success", s.conf.Hosts)
	return nil
}

// Disconnect 端口连接
func (s *Subscriber) Disconnect() error {
	s.cancel()

	if s.comsummer != nil {
		if err := s.comsummer.Close(); err != nil {
			s.l.Errorf("Failed to close consumer: ", err)
		}
	}

	return nil
}

// Sub 订阅事件
func (s *Subscriber) Sub(topic string, h bus.EventHandler) error {
	if topic == "" {
		return fmt.Errorf("topic required")
	}
	ctx, cancel := context.WithCancel(context.Background())
	s.h = h
	s.cancel = cancel
	go func() {
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := s.comsummer.Consume(ctx, strings.Split(topic, ","), s); err != nil {
				s.l.Errorf("Error from consumer: %v", err)
			}

			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			s.ready = make(chan bool)
		}
	}()

	<-s.ready // Await till the consumer has been set up
	s.l.Debugf("subscribe topics: %s, consumer group: %s ...", topic, s.conf.GroupID)
	return nil
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (s *Subscriber) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(s.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (s *Subscriber) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (s *Subscriber) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	// 具体消费消息
	for msg := range claim.Messages() {
		e := event.NewDefaultEvent()
		if err := json.Unmarshal(msg.Value, e); err != nil {
			s.l.Errorf("unmarshal data to event error, %s", err)
			continue
		}
		if s.h == nil {
			s.l.Error("event handler is nil")
			continue
		}
		s.h(msg.Topic, e)
		//run.Run(msg)
		// 更新位移
		session.MarkMessage(msg, "")
	}
	return nil
}
