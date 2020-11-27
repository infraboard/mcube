package kafka

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

func defaultBaseConfig() *baseConfig {
	return &baseConfig{
		Hosts:   []string{"127.0.0.1:9092"},
		Timeout: 30 * time.Second,
		Metadata: metaConfig{
			Retry: metaRetryConfig{
				Max:     5,
				Backoff: 300 * time.Millisecond,
			},
			RefreshFreq: 5 * time.Minute,
			Full:        true,
		},
		KeepAlive: 0,
		Version:   "1.0.0",

		ClientID:       "mcube",
		ChanBufferSize: 256,
		Username:       "",
		Password:       "",
	}
}

// Config 配置
type baseConfig struct {
	// 基本配置
	Hosts    []string `json:"hosts" yaml:"hosts" toml:"hosts" validate:"required" env:"BUS_KAFKA_HOSTS" envSeparator:","`
	Username string   `json:"username" yaml:"username" toml:"username" env:"BUS_KAFKA_USERNAME"`
	Password string   `json:"password" yaml:"password" toml:"password" env:"BUS_KAFKA_PASSWORD"`
	ClientID string   `json:"client_id" yaml:"client_id" toml:"client_id" env:"BUS_KAFKA_CLIENT_ID"`
	Version  string   `json:"version" yaml:"version" toml:"version" env:"BUS_KAFKA_VERSION"`

	// 优化配置
	KeepAlive      time.Duration `json:"keep_alive" yaml:"keep_alive" toml:"keep_alive" env:"BUS_KAFKA_KEEP_ALIVE"`
	Timeout        time.Duration `json:"timeout" yaml:"timeout" toml:"timeout" env:"BUS_KAFKA_TIMEOUT"`
	ChanBufferSize int           `json:"channel_buffer_size" yaml:"channel_buffer_size" toml:"channel_buffer_size" env:"BUS_KAFKA_CHANNEL_BUFFER_SIZE"`

	// meta配置和SSL配置
	Metadata metaConfig `json:"metadata" yaml:"metadata" toml:"metadata" env:"BUS_KAFKA_META_DATA"`
	Sasl     saslConfig `json:"sasl" yaml:"sasl" toml:"sasl"`
}

func (b *baseConfig) Validate() error {
	if len(b.Hosts) == 0 {
		return errors.New("no hosts configured")
	}

	if b.Username != "" && b.Password == "" {
		return fmt.Errorf("password must be set when username is configured")
	}

	_, err := sarama.ParseKafkaVersion(b.Version)
	if err != nil {
		return err
	}

	return nil
}

func (b *baseConfig) newBaseSaramaConfig() (*sarama.Config, error) {
	k := sarama.NewConfig()
	// configure network level properties
	timeout := b.Timeout
	k.Net.DialTimeout = timeout
	k.Net.ReadTimeout = timeout
	k.Net.WriteTimeout = timeout
	k.Net.KeepAlive = b.KeepAlive

	if b.Username != "" {
		k.Net.SASL.Enable = true
		k.Net.SASL.User = b.Username
		k.Net.SASL.Password = b.Password
		err := b.Sasl.configureSarama(k)
		if err != nil {
			return nil, err
		}
	}

	// configure metadata update properties
	k.Metadata.Retry.Max = b.Metadata.Retry.Max
	k.Metadata.Retry.Backoff = b.Metadata.Retry.Backoff
	k.Metadata.RefreshFrequency = b.Metadata.RefreshFreq
	k.Metadata.Full = b.Metadata.Full

	// configure per broker go channel buffering
	k.ChannelBufferSize = b.ChanBufferSize

	// configure client ID
	k.ClientID = b.ClientID

	version, err := sarama.ParseKafkaVersion(b.Version)
	if err != nil {
		return nil, fmt.Errorf("Unknown/unsupported kafka version: %s", err)
	}
	k.Version = version

	return k, nil
}

type metaConfig struct {
	Retry       metaRetryConfig `json:"retry" yaml:"retry" toml:"retry"`
	RefreshFreq time.Duration   `json:"refresh_frequency" yaml:"refresh_frequency" toml:"refresh_frequency" env:"BUS_KAFKA_META_REFRESH_REQ" validate:"min=0"`
	Full        bool            `json:"full" yaml:"full" toml:"full" env:"BUS_KAFKA_META_FULL"`
}

type metaRetryConfig struct {
	Max     int           `json:"max" yaml:"max" toml:"max" env:"BUS_KAFKA_META_RETRY_MAX" validate:"min=0"`
	Backoff time.Duration `json:"backoff" yaml:"backoff" toml:"backoff" env:"BUS_KAFKA_MATA_RETRY_BACKOFF" validate:"min=0"`
}

const (
	saslTypePlaintext   = sarama.SASLTypePlaintext
	saslTypeSCRAMSHA256 = sarama.SASLTypeSCRAMSHA256
	saslTypeSCRAMSHA512 = sarama.SASLTypeSCRAMSHA512
)

type saslConfig struct {
	SaslMechanism string `json:"mechanism" yaml:"mechanism" toml:"mechanism" env:"BUS_KAFKA_MECHANISM"`
}

func (c *saslConfig) configureSarama(config *sarama.Config) error {
	switch strings.ToUpper(c.SaslMechanism) { // try not to force users to use all upper case
	case "":
		// SASL is not enabled
		return nil
	case saslTypePlaintext:
		config.Net.SASL.Mechanism = sarama.SASLMechanism(sarama.SASLTypePlaintext)
	case saslTypeSCRAMSHA256:
		config.Net.SASL.Handshake = true
		config.Net.SASL.Mechanism = sarama.SASLMechanism(sarama.SASLTypeSCRAMSHA256)
		config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
			return &XDGSCRAMClient{HashGeneratorFcn: SHA256}
		}
	case saslTypeSCRAMSHA512:
		config.Net.SASL.Handshake = true
		config.Net.SASL.Mechanism = sarama.SASLMechanism(sarama.SASLTypeSCRAMSHA512)
		config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
			return &XDGSCRAMClient{HashGeneratorFcn: SHA512}
		}
	default:
		return fmt.Errorf("not valid mechanism '%v', only supported with PLAIN|SCRAM-SHA-512|SCRAM-SHA-256", c.SaslMechanism)
	}

	return nil
}
