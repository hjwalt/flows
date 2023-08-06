package flows

import (
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/runtime"
)

// Wiring configuration
type StatelessSingleFunctionConfiguration struct {
	KafkaProducerConfiguration []runtime.Configuration[*runtime_sarama.Producer]
	KafkaConsumerConfiguration []runtime.Configuration[*runtime_sarama.Consumer]
	RetryConfiguration         []runtime.Configuration[*runtime_retry.Retry]
	StatelessFunction          stateless.SingleFunction
	RouteConfiguration         []runtime.Configuration[*runtime_bunrouter.Router]
	AdditionalRuntimes         []runtime.Runtime
}

func (c StatelessSingleFunctionConfiguration) Runtime() runtime.Runtime {
	// producer runtime
	producer := KafkaProducer(c.KafkaProducerConfiguration)

	// function wrapping
	// - produce output messages
	messagesProduced := WrapSingleProduce(c.StatelessFunction, producer)

	// - retry
	produceRetry, retryRuntime := WrapRetry(messagesProduced, c.RetryConfiguration)

	// consumer runtime
	consumer := KafkaConsumerSingle(produceRetry, c.KafkaConsumerConfiguration)

	// http runtime
	routerRuntime := RouteRuntime(producer, c.RouteConfiguration)

	// add additional runtimes
	runtimes := []runtime.Runtime{
		routerRuntime,
		producer,
		consumer,
		retryRuntime,
	}
	if len(c.AdditionalRuntimes) > 0 {
		runtimes = append(c.AdditionalRuntimes, runtimes...)
	}

	return &RuntimeFacade{
		Runtimes: runtimes,
	}
}
