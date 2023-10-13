package kafka

import (
	"github.com/infraboard/mcube/ioc"
	kafka "github.com/segmentio/kafka-go"
)

const (
	KAFKA = "kafka"
)

type ScramAlgorithm string

const (
	SHA256 ScramAlgorithm = "SHA256"
	SHA512 ScramAlgorithm = "SHA512"
)

func Producer(topic string) *kafka.Writer {
	return ioc.Config().Get(KAFKA).(*Kafka).Producer(topic)
}

func ConsumerGroup(groupId string, topics []string) *kafka.Reader {
	return ioc.Config().Get(KAFKA).(*Kafka).ConsumerGroup(groupId, topics)
}
