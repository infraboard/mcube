package kafka

import (
	"github.com/infraboard/mcube/ioc"
	kafka "github.com/segmentio/kafka-go"
)

func init() {
	ioc.Config().Registry(&Kafka{
		Brokers:  []string{"127.0.0.1:9092"},
		UserName: "root",
		Password: "123456",
	})
}

type Kafka struct {
	Brokers  []string `toml:"brokers" json:"brokers" yaml:"brokers"  env:"KAFKA_BROKERS"`
	UserName string   `toml:"username" json:"username" yaml:"username"  env:"KAFKA_USERNAME"`
	Password string   `toml:"password" json:"password" yaml:"password"  env:"KAFKA_PASSWORD"`

	ioc.ObjectImpl
}

func (m *Kafka) Name() string {
	return KAFKA
}

func (m *Kafka) Init() error {
	return nil
}

func (k *Kafka) Producer(topic string) *kafka.Writer {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(k.Brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}
	return w
}

func (k *Kafka) Consumer() *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  k.Brokers,
		MaxBytes: 10e6, // 10MB
	})
}
