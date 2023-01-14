package runtime_sarama

import (
	"github.com/Shopify/sarama"
	"github.com/hjwalt/flows/runtime"
)

func WithConsumerRuntimeController(controller runtime.Controller) runtime.Configuration[*Consumer] {
	return func(c *Consumer) *Consumer {
		c.Controller = controller
		return c
	}
}

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

func WithConsumerSaramaConfig(saramaConfig *sarama.Config) runtime.Configuration[*Consumer] {
	return func(c *Consumer) *Consumer {
		c.SaramaConfiguration = saramaConfig
		return c
	}
}

func WithConsumerLoop(loop ConsumerLoop) runtime.Configuration[*Consumer] {
	return func(c *Consumer) *Consumer {
		c.Loop = loop
		return c
	}
}
