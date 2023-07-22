package runtime_sarama

import (
	"github.com/hjwalt/runway/runtime"
)

func WithConsumerTopic(topics ...string) runtime.Configuration[*Consumer] {
	return func(c *Consumer) *Consumer {
		c.Topics = topics
		return c
	}
}

func WithConsumerBroker(brokers ...string) runtime.Configuration[*Consumer] {
	return func(c *Consumer) *Consumer {
		c.Brokers = brokers
		return c
	}
}

func WithConsumerGroupName(groupName string) runtime.Configuration[*Consumer] {
	return func(c *Consumer) *Consumer {
		c.GroupName = groupName
		return c
	}
}

func WithConsumerSaramaConfigModifier(modifier SaramaConfigModifier) runtime.Configuration[*Consumer] {
	return func(c *Consumer) *Consumer {
		modifier(c.SaramaConfiguration)
		return c
	}
}

func WithConsumerLoop(loop ConsumerLoop) runtime.Configuration[*Consumer] {
	return func(c *Consumer) *Consumer {
		c.Handler = loop
		return c
	}
}
