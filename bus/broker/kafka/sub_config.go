package kafka

import (
	"fmt"
	"strings"

	"github.com/Shopify/sarama"
)

var balanceStrategyModes = map[string]sarama.BalanceStrategy{
	"":           sarama.BalanceStrategyRange,
	"sticky":     sarama.BalanceStrategySticky,
	"roundrobin": sarama.BalanceStrategyRoundRobin,
	"range":      sarama.BalanceStrategyRange,
}

// defaultSubscriberConfig 默认配置
func defaultSubscriberConfig() *subscriberConfig {
	return &subscriberConfig{
		GroupID:         "default",
		Offset:          "newest",
		BalanceStrategy: "range",
	}
}

// subscriberConfig todo
type subscriberConfig struct {
	GroupID         string `json:"group_id" yaml:"group_id" toml:"group_id" env:"BUS_KAFKA_SUBSCRIBER_GROUP_ID"`
	Offset          string `json:"offset" yaml:"offset" toml:"offset" env:"BUS_KAFKA_SUBSCRIBER_OFFSET"`
	BalanceStrategy string `json:"balance_strategy" yaml:"balance_strategy" toml:"balance_strategy" env:"BUS_KAFKA_SUBSCRIBER_BALANCE_STRATEGY"`
}

// Validate 校验配置
func (s *subscriberConfig) Validate() error {
	if err := validate.Struct(s); err != nil {
		return err
	}

	if _, ok := balanceStrategyModes[strings.ToLower(s.BalanceStrategy)]; !ok {
		return fmt.Errorf("balance_strategy mode '%v' unknown", s.BalanceStrategy)
	}

	return nil
}

func newSaramaSubConfig(base *baseConfig, conf *subscriberConfig) (*sarama.Config, error) {
	k, err := base.newBaseSaramaConfig()
	if err != nil {
		return nil, err
	}

	switch conf.Offset {
	case "oldest":
		k.Consumer.Offsets.Initial = sarama.OffsetOldest
	case "newest":
		k.Consumer.Offsets.Initial = sarama.OffsetNewest
	default:
		return nil, fmt.Errorf("-offset should be `oldest` or `newest`")
	}

	balanceStrategy, ok := balanceStrategyModes[strings.ToLower(conf.BalanceStrategy)]
	if !ok {
		return nil, fmt.Errorf("Unknown balance_strategy mode: '%v'", conf.BalanceStrategy)
	}
	k.Consumer.Group.Rebalance.Strategy = balanceStrategy

	if err := k.Validate(); err != nil {
		return nil, fmt.Errorf("Invalid kafka configuration: %v", err)
	}
	return k, nil
}
