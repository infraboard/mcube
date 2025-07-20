package nats_test

import (
	"testing"
	"time"

	"github.com/infraboard/mcube/v2/ioc"
	ioc_nats "github.com/infraboard/mcube/v2/ioc/config/nats"
	"github.com/nats-io/nats.go"
)

const (
	TEST_SUBJECT = "event_bus"
)

func TestPublish(t *testing.T) {
	err := ioc_nats.Get().Publish(TEST_SUBJECT, []byte("test"))
	if err != nil {
		t.Fatal(err)
	}

	err = ioc_nats.Get().Flush()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSub(t *testing.T) {
	_, err := ioc_nats.Get().Subscribe(TEST_SUBJECT, func(msg *nats.Msg) {
		t.Log(string(msg.Data))
	})
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(30 * time.Second)
}

func init() {
	err := ioc.ConfigIocObject(ioc.NewLoadConfigRequest())
	if err != nil {
		panic(err)
	}
}
