package kafka_test

import (
	"testing"

	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/ioc/config/kafka"
)

func TestKafkaProducer(t *testing.T) {
	m := kafka.Producer("test")
	t.Log(m)
}

func TestKafkaConsumerGroup(t *testing.T) {
	m := kafka.ConsumerGroup("test", []string{"test"})
	t.Log(m)
}

func init() {
	err := ioc.ConfigIocObject(ioc.NewLoadConfigRequest())
	if err != nil {
		panic(err)
	}
}
