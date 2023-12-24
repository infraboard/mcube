package kafka_test

import (
	"testing"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/kafka"
	"github.com/infraboard/mcube/v2/tools/file"
)

func TestKafkaProducer(t *testing.T) {
	m := kafka.Producer("test")
	t.Log(m)
}

func TestKafkaConsumerGroup(t *testing.T) {
	m := kafka.ConsumerGroup("test", []string{"test"})
	t.Log(m)
}

func TestDefaultConfig(t *testing.T) {
	file.MustToToml(
		kafka.AppName,
		ioc.Config().Get(kafka.AppName),
		"test/default.toml",
	)
}

func init() {
	err := ioc.ConfigIocObject(ioc.NewLoadConfigRequest())
	if err != nil {
		panic(err)
	}
}
