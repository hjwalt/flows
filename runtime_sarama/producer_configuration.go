package runtime_sarama

import (
	"github.com/hjwalt/runway/runtime"
)

func WithProducerBroker(brokers ...string) runtime.Configuration[*Producer] {
	return func(p *Producer) *Producer {
		p.Brokers = brokers
		return p
	}
}

func WithProducerSaramaConfigModifier(modifier SaramaConfigModifier) runtime.Configuration[*Producer] {
	return func(c *Producer) *Producer {
		modifier(c.SaramaConfiguration)
		return c
	}
}
