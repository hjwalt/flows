package flows

import (
	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateless"
)

// Wiring configuration
type StatelessSingleFunctionConfiguration struct {
	KafkaProducerConfiguration []runtime.Configuration[*runtime_sarama.Producer]
	KafkaConsumerConfiguration []runtime.Configuration[*runtime_sarama.Consumer]
	StatelessFunction          stateless.SingleFunction
	RouteConfiguration         []runtime.Configuration[*runtime_bunrouter.Router]
}

func (c StatelessSingleFunctionConfiguration) Runtime() runtime.Runtime {

	ctrl := runtime.NewController()

	// producer runtime
	producer := KafkaProducer(ctrl, c.KafkaProducerConfiguration)

	// function wrapping
	// - produce output messages
	messagesProduced := WrapSingleProduce(c.StatelessFunction, producer)

	// - retry
	produceRetry, retryRuntime := WrapRetry(messagesProduced)

	// consumer runtime
	consumer := KafkaConsumerSingle(ctrl, produceRetry, c.KafkaConsumerConfiguration)

	// http runtime
	routerRuntime := RouteRuntime(producer, c.RouteConfiguration)

	// multi runtime configuration
	multi := runtime.NewMulti(
		runtime.WithController(ctrl),
		runtime.WithRuntime(routerRuntime),
		runtime.WithRuntime(producer),
		runtime.WithRuntime(consumer),
		runtime.WithRuntime(retryRuntime),
	)
	return multi
}
