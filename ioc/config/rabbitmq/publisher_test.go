package rabbitmq_test

import (
	"testing"

	"github.com/infraboard/mcube/v2/ioc/config/rabbitmq"
)

const (
	TEST_SUBJECT = "event_bus"
)

func TestPublish(t *testing.T) {
	publisher, err := rabbitmq.NewPublisher()
	if err != nil {
		t.Fatal(err)
	}
	defer publisher.Close()

	msg := rabbitmq.NewFanoutMessage(TEST_SUBJECT, []byte("test"))
	err = publisher.Publish(t.Context(), msg)
	if err != nil {
		t.Fatal(err)
	}
}
