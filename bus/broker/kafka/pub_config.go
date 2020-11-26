package kafka

import (
	"fmt"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
)

var compressionModes = map[string]sarama.CompressionCodec{
	"none":   sarama.CompressionNone,
	"no":     sarama.CompressionNone,
	"off":    sarama.CompressionNone,
	"gzip":   sarama.CompressionGZIP,
	"lz4":    sarama.CompressionLZ4,
	"snappy": sarama.CompressionSnappy,
	"zstd":   sarama.CompressionZSTD,
}

var partitionerModes = map[string]sarama.PartitionerConstructor{
	"":            sarama.NewHashPartitioner,
	"hash":        sarama.NewHashPartitioner,
	"random":      sarama.NewRandomPartitioner,
	"manual":      sarama.NewManualPartitioner,
	"round_robin": sarama.NewRoundRobinPartitioner,
}

// defaultPublisherConfig 默认配置
func defaultPublisherConfig() *publisherConfig {
	return &publisherConfig{
		BulkMaxSize:        2048,
		BulkFlushFrequency: 0,
		MaxMessageBytes:    nil, // use library default
		RequiredACKs:       nil, // use library default
		Compression:        "gzip",
		CompressionLevel:   4,
		MaxRetries:         3,
	}
}

// PublisherConfig todo
type publisherConfig struct {
	BulkMaxSize        int           `json:"bulk_max_size" yaml:"bulk_max_size" toml:"bulk_max_size" env:"BUS_KAFKA_PUBLISHER_BULK_MAX_SIZE"`
	PublishTimeout     time.Duration `json:"publish_timeout" yaml:"publish_timeout" toml:"publish_timeout" env:"BUS_KAFKA_PUBLISHER_TIMEOUT"`
	BulkFlushFrequency time.Duration `json:"bulk_flush_frequency" yaml:"bulk_flush_frequency" toml:"bulk_flush_frequency" env:"BUS_KAFKA_PUBLISHER_BULK_FLUSH_FREQUENCY"`
	Partitioner        string        `json:"partitioner" yaml:"partitioner" toml:"partitioner" env:"BUS_KAFKA_PUBLISHER_PARTITIONER"`
	Compression        string        `json:"compression" yaml:"compression" toml:"compression" env:"BUS_KAFKA_PUBLISHER_COMPRESSION"`
	CompressionLevel   int           `json:"compression_level" yaml:"compression_level" toml:"compression_level" env:"BUS_KAFKA_PUBLISHER_COMPRESSION_LEVEL"`
	MaxRetries         int           `json:"max_retries" yaml:"max_retries" toml:"max_retries" env:"BUS_KAFKA_PUBLISHER_MAX_RETRIES"`
	MaxMessageBytes    *int          `json:"max_message_bytes" yaml:"max_message_bytes" toml:"max_message_bytes" env:"BUS_KAFKA_PUBLISHER_MAX_MESSAGE_BYTES"`
	RequiredACKs       *int          `json:"required_acks" yaml:"required_acks" toml:"required_acks" env:"BUS_KAFKA_PUBLISHER_REQUIRED_ACKS"`
}

// Validate 校验配置
func (c *publisherConfig) Validate() error {
	if err := validate.Struct(c); err != nil {
		return err
	}

	if _, ok := compressionModes[strings.ToLower(c.Compression)]; !ok {
		return fmt.Errorf("compression mode '%v' unknown", c.Compression)
	}

	if _, ok := partitionerModes[strings.ToLower(c.Partitioner)]; !ok {
		return fmt.Errorf("partitioner mode '%v' unknown", c.Partitioner)
	}

	if c.Compression == "gzip" {
		lvl := c.CompressionLevel
		if lvl != sarama.CompressionLevelDefault && !(0 <= lvl && lvl <= 9) {
			return fmt.Errorf("compression_level must be between 0 and 9")
		}
	}
	return nil
}

func newSaramaPubConfig(base *baseConfig, config *publisherConfig) (*sarama.Config, error) {
	k, err := base.newBaseSaramaConfig()
	if err != nil {
		return nil, err
	}

	// configure producer API properties
	if config.MaxMessageBytes != nil {
		k.Producer.MaxMessageBytes = *config.MaxMessageBytes
	}
	if config.RequiredACKs != nil {
		k.Producer.RequiredAcks = sarama.RequiredAcks(*config.RequiredACKs)
	}

	// 压缩模式
	compressionMode, ok := compressionModes[strings.ToLower(config.Compression)]
	if !ok {
		return nil, fmt.Errorf("Unknown compression mode: '%v'", config.Compression)
	}
	k.Producer.Compression = compressionMode

	// partion 模式
	partitionerMode, ok := partitionerModes[strings.ToLower(config.Partitioner)]
	if !ok {
		return nil, fmt.Errorf("Unknown partition mode: '%v'", config.Partitioner)
	}
	k.Producer.Partitioner = partitionerMode

	k.Producer.Return.Successes = true // enable return channel for signaling
	k.Producer.Return.Errors = true

	retryMax := config.MaxRetries
	if retryMax < 0 {
		retryMax = 1000
	}
	k.Producer.Retry.Max = retryMax

	// configure bulk size
	k.Producer.Flush.MaxMessages = config.BulkMaxSize
	if config.BulkFlushFrequency > 0 {
		k.Producer.Flush.Frequency = config.BulkFlushFrequency
	}

	if err := k.Validate(); err != nil {
		return nil, fmt.Errorf("Invalid kafka configuration: %v", err)
	}
	return k, nil
}
