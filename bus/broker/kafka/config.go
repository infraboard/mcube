package kafka

// NewDefultConfig todo
func NewDefultConfig() *Config {
	return &Config{
		baseConfig:       defaultBaseConfig(),
		publisherConfig:  defaultPublisherConfig(),
		subscriberConfig: defaultSubscriberConfig(),
	}
}

// Config 配置
type Config struct {
	*baseConfig       `yaml:",inline"`
	*publisherConfig  `yaml:",inline"`
	*subscriberConfig `yaml:",inline"`
}

// ValidatePublisherConfig todo
func (conf *Config) ValidatePublisherConfig() error {
	if err := conf.baseConfig.Validate(); err != nil {
		return err
	}
	return conf.publisherConfig.Validate()
}

// ValidateSubscriberConfig todo
func (conf *Config) ValidateSubscriberConfig() error {
	if err := conf.baseConfig.Validate(); err != nil {
		return err
	}
	return conf.subscriberConfig.Validate()
}
