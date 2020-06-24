package nats_test

import (
	"testing"

	"github.com/infraboard/mcube/bus/broker/nats"
	"github.com/infraboard/mcube/bus/event"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/stretchr/testify/assert"
)

func TestPubSub(t *testing.T) {
	should := assert.New(t)
	b := nats.NewBroker(nats.NewDefaultConfig())
	b.Debug(zap.L().Named("Nats Bus"))

	sourceEvent := event.NewEvent()

	should.NoError(b.Connect())

	err := b.Sub("test", func(e *event.Event) error {
		should.Equal(sourceEvent.ID, e.ID)
		return nil
	})
	should.NoError(err)

	should.NoError(b.Pub("test", sourceEvent))
	should.NoError(b.Disconnect())
}

func init() {
	zap.DevelopmentSetup()
}
