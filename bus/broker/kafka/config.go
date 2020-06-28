package kafka

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
)

// NewDefaultConfig 默认配置
func NewDefaultConfig() *Config {
	return &Config{
		Hosts:              nil,
		Timeout:            30 * time.Second,
		BulkMaxSize:        2048,
		BulkFlushFrequency: 0,
		Metadata: metaConfig{
			Retry: metaRetryConfig{
				Max:     3,
				Backoff: 250 * time.Millisecond,
			},
			RefreshFreq: 10 * time.Minute,
			Full:        false,
		},
		KeepAlive:        0,
		MaxMessageBytes:  nil, // use library default
		RequiredACKs:     nil, // use library default
		BrokerTimeout:    10 * time.Second,
		Compression:      "gzip",
		CompressionLevel: 4,
		Version:          "1.0.0",
		MaxRetries:       3,
		ClientID:         "mcube",
		ChanBufferSize:   256,
		Username:         "",
		Password:         "",
	}
}

// Config 配置
type Config struct {
	Hosts              []string      `json:"hosts"               validate:"required"`
	Timeout            time.Duration `json:"timeout"             validate:"min=1"`
	Metadata           metaConfig    `json:"metadata"`
	KeepAlive          time.Duration `json:"keep_alive"          validate:"min=0"`
	MaxMessageBytes    *int          `json:"max_message_bytes"   validate:"min=1"`
	RequiredACKs       *int          `json:"required_acks"       validate:"min=-1"`
	BrokerTimeout      time.Duration `json:"broker_timeout"      validate:"min=1"`
	Compression        string        `json:"compression"`
	CompressionLevel   int           `json:"compression_level"`
	Version            string        `json:"version"`
	BulkMaxSize        int           `json:"bulk_max_size"`
	BulkFlushFrequency time.Duration `json:"bulk_flush_frequency"`
	MaxRetries         int           `json:"max_retries"         validate:"min=-1,nonzero"`
	ClientID           string        `json:"client_id"`
	ChanBufferSize     int           `json:"channel_buffer_size" validate:"min=1"`
	Username           string        `json:"username"`
	Password           string        `json:"password"`
	Sasl               saslConfig    `json:"sasl"`
}

// KafkaVersion todo
func (c *Config) KafkaVersion() sarama.KafkaVersion {
	v, _ := sarama.ParseKafkaVersion(c.Version)
	return v
}

// Validate 校验配置
func (c *Config) Validate() error {
	if err := validate.Struct(c); err != nil {
		return err
	}

	if len(c.Hosts) == 0 {
		return errors.New("no hosts configured")
	}

	if _, ok := compressionModes[strings.ToLower(c.Compression)]; !ok {
		return fmt.Errorf("compression mode '%v' unknown", c.Compression)
	}

	if _, err := sarama.ParseKafkaVersion(c.Version); err != nil {
		return err
	}

	if c.Username != "" && c.Password == "" {
		return fmt.Errorf("password must be set when username is configured")
	}

	if c.Compression == "gzip" {
		lvl := c.CompressionLevel
		if lvl != sarama.CompressionLevelDefault && !(0 <= lvl && lvl <= 9) {
			return fmt.Errorf("compression_level must be between 0 and 9")
		}
	}
	return nil
}

type saslConfig struct {
	SaslMechanism string `json:"mechanism"`
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

type metaConfig struct {
	Retry       metaRetryConfig `json:"retry"`
	RefreshFreq time.Duration   `json:"refresh_frequency" validate:"min=0"`
	Full        bool            `json:"full"`
}

type metaRetryConfig struct {
	Max     int           `json:"max"     validate:"min=0"`
	Backoff time.Duration `json:"backoff" validate:"min=0"`
}

var compressionModes = map[string]sarama.CompressionCodec{
	"none":   sarama.CompressionNone,
	"no":     sarama.CompressionNone,
	"off":    sarama.CompressionNone,
	"gzip":   sarama.CompressionGZIP,
	"lz4":    sarama.CompressionLZ4,
	"snappy": sarama.CompressionSnappy,
	"zstd":   sarama.CompressionZSTD,
}

const (
	saslTypePlaintext   = sarama.SASLTypePlaintext
	saslTypeSCRAMSHA256 = sarama.SASLTypeSCRAMSHA256
	saslTypeSCRAMSHA512 = sarama.SASLTypeSCRAMSHA512
)

func newSaramaConfig(config *Config) (*sarama.Config, error) {
	k := sarama.NewConfig()

	// configure network level properties
	timeout := config.Timeout
	k.Net.DialTimeout = timeout
	k.Net.ReadTimeout = timeout
	k.Net.WriteTimeout = timeout
	k.Net.KeepAlive = config.KeepAlive
	k.Producer.Timeout = config.BrokerTimeout
	k.Producer.CompressionLevel = config.CompressionLevel

	// tls, err := outputs.LoadTLSConfig(config.TLS)
	// if err != nil {
	// 	return nil, err
	// }

	// if tls != nil {
	// 	k.Net.TLS.Enable = true
	// 	k.Net.TLS.Config = tls.BuildModuleConfig("")
	// }

	if config.Username != "" {
		k.Net.SASL.Enable = true
		k.Net.SASL.User = config.Username
		k.Net.SASL.Password = config.Password
		err := config.Sasl.configureSarama(k)
		if err != nil {
			return nil, err
		}
	}

	// configure metadata update properties
	k.Metadata.Retry.Max = config.Metadata.Retry.Max
	k.Metadata.Retry.Backoff = config.Metadata.Retry.Backoff
	k.Metadata.RefreshFrequency = config.Metadata.RefreshFreq
	k.Metadata.Full = config.Metadata.Full

	// configure producer API properties
	if config.MaxMessageBytes != nil {
		k.Producer.MaxMessageBytes = *config.MaxMessageBytes
	}
	if config.RequiredACKs != nil {
		k.Producer.RequiredAcks = sarama.RequiredAcks(*config.RequiredACKs)
	}

	compressionMode, ok := compressionModes[strings.ToLower(config.Compression)]
	if !ok {
		return nil, fmt.Errorf("Unknown compression mode: '%v'", config.Compression)
	}
	k.Producer.Compression = compressionMode

	k.Producer.Return.Successes = true // enable return channel for signaling
	k.Producer.Return.Errors = true

	// have retries being handled by libbeat, disable retries in sarama library
	retryMax := config.MaxRetries
	if retryMax < 0 {
		retryMax = 1000
	}
	k.Producer.Retry.Max = retryMax
	// TODO: k.Producer.Retry.Backoff = ?

	// configure per broker go channel buffering
	k.ChannelBufferSize = config.ChanBufferSize

	// configure bulk size
	k.Producer.Flush.MaxMessages = config.BulkMaxSize
	if config.BulkFlushFrequency > 0 {
		k.Producer.Flush.Frequency = config.BulkFlushFrequency
	}

	// configure client ID
	k.ClientID = config.ClientID

	version, err := sarama.ParseKafkaVersion(config.Version)
	if err != nil {
		return nil, fmt.Errorf("Unknown/unsupported kafka version: %s", err)
	}
	k.Version = version

	if err := k.Validate(); err != nil {
		return nil, fmt.Errorf("Invalid kafka configuration: %v", err)
	}
	return k, nil
}
