package flows

import (
	"context"

	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
)

// Wiring configuration
type StatelessSingleFunctionConfiguration struct {
	KafkaProducerConfiguration []runtime.Configuration[*runtime_sarama.Producer]
	KafkaConsumerConfiguration []runtime.Configuration[*runtime_sarama.Consumer]
	RetryConfiguration         []runtime.Configuration[*runtime_retry.Retry]
	StatelessFunction          stateless.SingleFunction
	RouteConfiguration         []runtime.Configuration[*runtime_bunrouter.Router]
}

func (c StatelessSingleFunctionConfiguration) Runtime() runtime.Runtime {
	RegisterRetry(c.RetryConfiguration)
	RegisterProducerConfig(c.KafkaProducerConfiguration)
	RegisterProducer()
	RegisterConsumerSingleConfig(c.KafkaConsumerConfiguration)
	RegisterConsumerSingle()
	RegisterRoute(c.RouteConfiguration)
	inverse.Register[stateless.SingleFunction](QualifierKafkaConsumerSingleFunction, func(ctx context.Context) (stateless.SingleFunction, error) {
		retry, err := GetRetry(ctx)
		if err != nil {
			return nil, err
		}
		producer, err := GetKafkaProducer(ctx)
		if err != nil {
			return nil, err
		}

		wrappedFunction := c.StatelessFunction
		wrappedFunction = stateless.NewSingleProducer(
			stateless.WithSingleProducerNextFunction(wrappedFunction),
			stateless.WithSingleProducerRuntime(producer),
			stateless.WithSingleProducerPrometheus(),
		)
		wrappedFunction = stateless.NewSingleRetry(
			stateless.WithSingleRetryNextFunction(wrappedFunction),
			stateless.WithSingleRetryRuntime(retry),
			stateless.WithSingleRetryPrometheus(),
		)

		return wrappedFunction, nil
	})

	return &RuntimeFacade{
		Runtimes: InjectedRuntimes(),
	}
}
