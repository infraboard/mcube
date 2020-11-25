package kafka_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/infraboard/mcube/bus"
	"github.com/infraboard/mcube/bus/broker/kafka"
	"github.com/infraboard/mcube/bus/event"
	"github.com/infraboard/mcube/logger/zap"
)

var (
	sourceEvent = event.NewOperateEvent(nil, &event.OperateEvent{})
)

func TestPub(t *testing.T) {
	should := assert.New(t)
	conf := kafka.DefaultPublisherConfig()
	conf.Hosts = []string{"127.0.0.1:9092"}
	kp, err := kafka.NewPublisher(conf)
	bus.SetPublisher(kp)
	if should.NoError(err) {
		should.NoError(kp.Connect())
		should.NoError(kp.Pub("t1", sourceEvent))
	}

	time.Sleep(5 * time.Second)
}

func TestSub(t *testing.T) {
	should := assert.New(t)
	conf := kafka.DefaultSubscriberConfig()
	conf.Hosts = []string{"127.0.0.1:9092"}
	conf.GroupID = "sub-test"
	bus, err := kafka.NewSubscriber(conf)
	if should.NoError(err) {
		should.NoError(bus.Connect())
		bus.Sub("t1", func(topic string, e *event.Event) error {
			should.Equal("t1", topic)
			should.Equal(sourceEvent.ID, e.ID)
			t.Log("sub event: ", e)
			return nil
		})
	}

	time.Sleep(2 * time.Second)
}

func init() {
	zap.DevelopmentSetup()
}
