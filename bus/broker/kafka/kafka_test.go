package kafka_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/infraboard/mcube/bus/broker/kafka"
	"github.com/infraboard/mcube/bus/event"
	"github.com/infraboard/mcube/logger/zap"
)

func TestPub(t *testing.T) {
	sourceEvent := event.NewEvent()

	should := assert.New(t)
	conf := kafka.DefaultPublisherConfig()
	conf.Hosts = []string{"127.0.0.1:32774"}
	bus, err := kafka.NewPublisher(conf)
	if should.NoError(err) {
		should.NoError(bus.Connect())

		bus.Pub("t1", sourceEvent)
	}

	time.Sleep(2 * time.Second)
}

func TestSub(t *testing.T) {
	sourceEvent := event.NewEvent()

	should := assert.New(t)
	conf := kafka.DefaultSubscriberConfig()
	conf.Hosts = []string{"127.0.0.1:32774"}
	conf.GroupID = "sub-test"
	conf.Offset = "oldest"
	bus, err := kafka.NewSubscriber(conf)
	if should.NoError(err) {
		should.NoError(bus.Connect())

		bus.Sub("t1", func(e *event.Event) error {
			should.Equal(sourceEvent.ID, e.ID)
			return nil
		})
	}

	time.Sleep(2 * time.Second)
}

func init() {
	zap.DevelopmentSetup()
}
