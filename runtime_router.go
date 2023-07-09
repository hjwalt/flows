package flows

import (
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/runway/runtime"
)

// Wiring configuration
type RouterConfiguration struct {
	KafkaProducerConfiguration []runtime.Configuration[*runtime_sarama.Producer]
	RouteConfiguration         []runtime.Configuration[*runtime_bunrouter.Router]
}

func (c RouterConfiguration) Runtime() runtime.Runtime {
	// producer runtime
	producer := KafkaProducer(c.KafkaProducerConfiguration)

	// http runtime
	routerRuntime := RouteRuntime(producer, c.RouteConfiguration)

	return &RuntimeFacade{
		Runtimes: []runtime.Runtime{
			routerRuntime,
			producer,
		},
	}
}
