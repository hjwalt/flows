package runtime_sarama

import (
	"github.com/Shopify/sarama"
	"github.com/hjwalt/runway/runtime"
)

func WithProducerBroker(brokers ...string) runtime.Configuration[*Producer] {
	return func(p *Producer) *Producer {
		p.Brokers = brokers
		return p
	}
}

func WithProducerSaramaConfig(saramaConfig *sarama.Config) runtime.Configuration[*Producer] {
	return func(p *Producer) *Producer {
		p.SaramaConfiguration = saramaConfig
		return p
	}
}
