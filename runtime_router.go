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
	AdditionalRuntimes         []runtime.Runtime
}

func (c RouterConfiguration) Runtime() runtime.Runtime {
	// producer runtime
	producer := KafkaProducer(c.KafkaProducerConfiguration)

	// http runtime
	routerRuntime := RouteRuntime(producer, c.RouteConfiguration)

	// add additional runtimes
	runtimes := []runtime.Runtime{
		routerRuntime,
		producer,
	}
	if len(c.AdditionalRuntimes) > 0 {
		runtimes = append(c.AdditionalRuntimes, runtimes...)
	}

	return &RuntimeFacade{
		Runtimes: runtimes,
	}
}
