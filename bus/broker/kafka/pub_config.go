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

// DefaultPublisherConfig 默认配置
func DefaultPublisherConfig() *PublisherConfig {
	return &PublisherConfig{
		baseConfig:         defaultBaseConfig(),
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
type PublisherConfig struct {
	*baseConfig
	BulkMaxSize        int           `json:"bulk_max_size"`
	PublishTimeout     time.Duration `json:"publish_timeout"      `
	BulkFlushFrequency time.Duration `json:"bulk_flush_frequency"`
	Partitioner        string        `json:"partitioner"`
	Compression        string        `json:"compression"`
	CompressionLevel   int           `json:"compression_level"`
	MaxRetries         int           `json:"max_retries"`
	MaxMessageBytes    *int          `json:"max_message_bytes"`
	RequiredACKs       *int          `json:"required_acks"`
}

// Validate 校验配置
func (c *PublisherConfig) Validate() error {
	if err := validate.Struct(c); err != nil {
		return err
	}

	if err := c.baseConfig.validate(); err != nil {
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

func newSaramaPubConfig(config *PublisherConfig) (*sarama.Config, error) {
	k, err := config.newBaseSaramaConfig()
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
