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
	return ioc.Config().Get(AppName).(*Kafka).Producer(topic)
}

func ConsumerGroup(groupId string, topics []string) *kafka.Reader {
	return ioc.Config().Get(AppName).(*Kafka).ConsumerGroup(groupId, topics)
}
