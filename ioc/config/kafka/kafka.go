package kafka

import (
	"fmt"
	"time"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	kafka "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/scram"
)

func init() {
	ioc.Config().Registry(defaultConfig)
}

var defaultConfig = &Kafka{
	Brokers:        []string{"127.0.0.1:9092"},
	ScramAlgorithm: SHA512,
	Debug:          false,
}

type Kafka struct {
	Brokers        []string       `toml:"brokers" json:"brokers" yaml:"brokers"  env:"BROKERS"`
	ScramAlgorithm ScramAlgorithm `toml:"scram_algorithm" json:"scram_algorithm" yaml:"scram_algorithm"  env:"SCRAM_ALGORITHM"`
	UserName       string         `toml:"username" json:"username" yaml:"username"  env:"USERNAME"`
	Password       string         `toml:"password" json:"password" yaml:"password"  env:"PASSWORD"`
	Debug          bool           `toml:"debug" json:"debug" yaml:"debug"  env:"DEBUG"`

	mechanism sasl.Mechanism
	ioc.ObjectImpl
}

func (m *Kafka) Name() string {
	return AppName
}

func (i *Kafka) Priority() int {
	return 695
}

func (k *Kafka) Init() error {
	if k.UserName != "" {
		mechanism, err := scram.Mechanism(k.scramAlgorithm(), k.UserName, k.Password)
		if err != nil {
			return err
		}
		k.mechanism = mechanism
	}

	return nil
}

func (k *Kafka) scramAlgorithm() scram.Algorithm {
	switch k.ScramAlgorithm {
	case SHA256:
		return scram.SHA256
	default:
		return scram.SHA512
	}
}

func (k *Kafka) Producer(topic string) *kafka.Writer {
	l := log.Sub(fmt.Sprintf("producer_%s", topic))

	w := &kafka.Writer{
		Addr:                   kafka.TCP(k.Brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
		ErrorLogger:            kafka.LoggerFunc(l.Error().Msgf),
		Transport:              &kafka.Transport{SASL: k.mechanism},
	}
	if k.Debug {
		w.Logger = kafka.LoggerFunc(l.Debug().Msgf)
	}
	return w
}

func (k *Kafka) ConsumerGroup(groupId string, topics []string) *kafka.Reader {
	l := log.Sub(fmt.Sprintf("consumer_group_%s", groupId))
	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: k.mechanism,
	}

	conf := kafka.ReaderConfig{
		Brokers:     k.Brokers,
		Dialer:      dialer,
		GroupID:     groupId,
		GroupTopics: topics,
		ErrorLogger: kafka.LoggerFunc(l.Error().Msgf),
	}

	if k.Debug {
		conf.Logger = kafka.LoggerFunc(l.Debug().Msgf)
	}

	return kafka.NewReader(conf)
}
