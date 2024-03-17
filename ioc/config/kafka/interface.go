package kafka

import (
	"github.com/infraboard/mcube/v2/ioc"
	kafka "github.com/segmentio/kafka-go"
)

const (
	AppName = "kafka"
)

type ScramAlgorithm string

const (
	SHA256 ScramAlgorithm = "SHA256"
	SHA512 ScramAlgorithm = "SHA512"
)

func Producer(topic string) *kafka.Writer {
	return Get().Producer(topic)
}

func ConsumerGroup(groupId string, topics []string) *kafka.Reader {
	return Get().ConsumerGroup(groupId, topics)
}

func Get() *Kafka {
	obj := ioc.Config().Get(AppName)
	if obj == nil {
		return defaultConfig
	}
	return obj.(*Kafka)
}
