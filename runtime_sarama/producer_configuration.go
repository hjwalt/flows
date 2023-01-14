package runtime_sarama

import (
	"github.com/Shopify/sarama"
	"github.com/hjwalt/flows/runtime"
)

func WithProducerRuntimeController(controller runtime.Controller) runtime.Configuration[*Producer] {
	return func(p *Producer) *Producer {
		p.Controller = controller
		return p
	}
}

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
