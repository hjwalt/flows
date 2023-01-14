package flows

import (
	"time"

	"github.com/avast/retry-go"
	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateless"
)

// Wiring configuration
type StatelessSingleFunctionConfiguration struct {
	KafkaProducerConfiguration []runtime.Configuration[*runtime_sarama.Producer]
	KafkaConsumerConfiguration []runtime.Configuration[*runtime_sarama.Consumer]
	StatelessFunction          stateless.StatelessBinarySingleFunction
	RouteConfiguration         []runtime.Configuration[*runtime_bunrouter.Router]
}

func (c StatelessSingleFunctionConfiguration) Runtime() runtime.Runtime {

	ctrl := runtime.NewController()

	// producer runtime
	producerConfig := append(
		c.KafkaProducerConfiguration,
		runtime_sarama.WithProducerRuntimeController(ctrl),
	)
	producer := runtime_sarama.NewProducer(producerConfig...)

	// function wrapping
	// - produce output messages
	messagesProduced := stateless.NewSingleProducer(
		stateless.WithSingleProducerNextFunction(c.StatelessFunction),
		stateless.WithSingleProducerRuntime(producer),
	)

	// - retry
	retryRuntime := runtime_retry.NewRetry(
		runtime_retry.WithRetryOption(
			retry.Attempts(1000000),
			retry.Delay(10*time.Millisecond),
			retry.MaxDelay(time.Second),
			retry.MaxJitter(time.Second),
			retry.DelayType(retry.BackOffDelay),
		),
	)
	produceRetry := stateless.NewSingleRetry(
		stateless.WithSingleRetryRuntime(retryRuntime),
		stateless.WithSingleRetryNextFunction(messagesProduced),
	)

	// sarama consumer loop
	consumerLoop := runtime_sarama.NewSingleLoop(
		runtime_sarama.WithLoopSingleFunction(produceRetry),
	)

	// consumer runtime
	consumerConfig := append(
		c.KafkaConsumerConfiguration,
		runtime_sarama.WithConsumerRuntimeController(ctrl),
		runtime_sarama.WithConsumerLoop(consumerLoop),
	)
	consumer := runtime_sarama.NewConsumer(consumerConfig...)

	// http runtime, prometheus first for hard prometheus path
	routeConfig := append(
		make([]runtime.Configuration[*runtime_bunrouter.Router], 0),
		runtime_bunrouter.WithRouterPrometheus(),
	)
	routeConfig = append(
		routeConfig,
		c.RouteConfiguration...,
	)
	routerRuntime := runtime_bunrouter.NewRouter(routeConfig...)

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
