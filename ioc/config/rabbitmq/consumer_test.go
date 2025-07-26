package rabbitmq_test

import (
	"context"
	"testing"
	"time"

	"github.com/infraboard/mcube/v2/ioc/config/rabbitmq"
)

const (
	FANOUT_SUBJECT = "event_bus_fanout"
	TOPIC_SUBJECT  = "event_bus_topic"
)

func TestFanoutSubscribePublish(t *testing.T) {
	publisher, err := rabbitmq.NewPublisher()
	if err != nil {
		t.Fatal(err)
	}
	defer publisher.Close()

	msg := rabbitmq.NewFanoutMessage(FANOUT_SUBJECT, []byte("test"))
	err = publisher.Publish(t.Context(), msg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFanoutSubscribe(t *testing.T) {
	consumer1, err := rabbitmq.NewConsumer()
	if err != nil {
		t.Fatal(err)
	}

	err = consumer1.FanoutSubscribe(t.Context(), FANOUT_SUBJECT, func(ctx context.Context, msg *rabbitmq.Message) error {
		t.Log("consumer1: " + string(msg.Body))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	consumer2, err := rabbitmq.NewConsumer()
	if err != nil {
		t.Fatal(err)
	}

	err = consumer2.FanoutSubscribe(t.Context(), FANOUT_SUBJECT, func(ctx context.Context, msg *rabbitmq.Message) error {
		t.Log("consumer2: " + string(msg.Body))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(60 * time.Second)
}

func TestTopicSubscribePublish(t *testing.T) {
	publisher, err := rabbitmq.NewPublisher()
	if err != nil {
		t.Fatal(err)
	}
	defer publisher.Close()

	logMsg := rabbitmq.NewTopicMessage(TOPIC_SUBJECT, "event_bus.logs.info", []byte("log"))
	err = publisher.Publish(t.Context(), logMsg)
	if err != nil {
		t.Fatal(err)
	}

	alertMsg := rabbitmq.NewTopicMessage(TOPIC_SUBJECT, "event_bus.alerts.info", []byte("alert"))
	err = publisher.Publish(t.Context(), alertMsg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTopicSubscribe(t *testing.T) {
	logConsumer, err := rabbitmq.NewConsumer()
	if err != nil {
		t.Fatal(err)
	}

	err = logConsumer.TopicSubscribe(t.Context(), TOPIC_SUBJECT, "event_bus.logs.*", func(ctx context.Context, msg *rabbitmq.Message) error {
		t.Log(string(msg.Body))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	alertConsumer, err := rabbitmq.NewConsumer()
	if err != nil {
		t.Fatal(err)
	}

	err = alertConsumer.TopicSubscribe(t.Context(), TOPIC_SUBJECT, "event_bus.alerts.*", func(ctx context.Context, msg *rabbitmq.Message) error {
		t.Log(string(msg.Body))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(60 * time.Second)
}

func TestDirectSubscribePublish(t *testing.T) {
	publisher, err := rabbitmq.NewPublisher()
	if err != nil {
		t.Fatal(err)
	}
	defer publisher.Close()

	msg := rabbitmq.NewQueueMessage("orders", []byte("test1"))
	err = publisher.Publish(t.Context(), msg)
	if err != nil {
		t.Fatal(err)
	}

	msg = rabbitmq.NewQueueMessage("orders", []byte("test2"))
	err = publisher.Publish(t.Context(), msg)
	if err != nil {
		t.Fatal(err)
	}

	msg = rabbitmq.NewQueueMessage("orders", []byte("test3"))
	err = publisher.Publish(t.Context(), msg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDirectSubscribe(t *testing.T) {
	consumer_a, err := rabbitmq.NewConsumer()
	if err != nil {
		t.Fatal(err)
	}

	err = consumer_a.DirectSubscribe(t.Context(), "orders", func(ctx context.Context, msg *rabbitmq.Message) error {
		t.Log("a: " + string(msg.Body))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	consumer_b, err := rabbitmq.NewConsumer()
	if err != nil {
		t.Fatal(err)
	}

	err = consumer_b.DirectSubscribe(t.Context(), "orders", func(ctx context.Context, msg *rabbitmq.Message) error {
		t.Log("b: " + string(msg.Body))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(60 * time.Second)
}
