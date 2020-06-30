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

// DefaultSubscriberConfig 默认配置
func DefaultSubscriberConfig() *SubscriberConfig {
	return &SubscriberConfig{
		baseConfig:      defaultBaseConfig(),
		GroupID:         "default",
		Offset:          "newest",
		BalanceStrategy: "range",
	}
}

// SubscriberConfig todo
type SubscriberConfig struct {
	*baseConfig
	GroupID         string `json:"group_id,omitempty"`
	Offset          string `json:"offset,omitempty"`
	BalanceStrategy string `json:"balance_strategy,omitempty"`
}

// Validate 校验配置
func (s *SubscriberConfig) Validate() error {
	if err := validate.Struct(s); err != nil {
		return err
	}

	if err := s.baseConfig.validate(); err != nil {
		return err
	}

	if _, ok := balanceStrategyModes[strings.ToLower(s.BalanceStrategy)]; !ok {
		return fmt.Errorf("balance_strategy mode '%v' unknown", s.BalanceStrategy)
	}

	return nil
}

func newSaramaSubConfig(conf *SubscriberConfig) (*sarama.Config, error) {
	k, err := conf.newBaseSaramaConfig()
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
