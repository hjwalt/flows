package flows

import (
	"context"

	"github.com/hjwalt/flows/collect"
	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
)

// Wiring configuration
type CollectorConfiguration struct {
	KafkaProducerConfiguration []runtime.Configuration[*runtime_sarama.Producer]
	KafkaConsumerConfiguration []runtime.Configuration[*runtime_sarama.Consumer]
	RetryConfiguration         []runtime.Configuration[*runtime_retry.Retry]
	PersistenceIdFunction      stateful.PersistenceIdFunction[[]byte, []byte]
	Aggregator                 collect.Aggregator[message.Bytes, message.Bytes, message.Bytes]
	Collector                  collect.Collector
	RouteConfiguration         []runtime.Configuration[*runtime_bunrouter.Router]
}

func (c CollectorConfiguration) Register() {
	RegisterRetry(c.RetryConfiguration)
	RegisterProducerConfig(c.KafkaProducerConfiguration)
	RegisterProducer()
	RegisterConsumerKeyedConfig(c.KafkaConsumerConfiguration)
	RegisterConsumer()
	RegisterRoute(c.RouteConfiguration)
	inverse.RegisterInstance[stateful.PersistenceIdFunction[message.Bytes, message.Bytes]](QualifierKafkaConsumerKeyFunction, c.PersistenceIdFunction)
	inverse.Register[stateless.BatchFunction](QualifierKafkaConsumerBatchFunction, func(ctx context.Context) (stateless.BatchFunction, error) {
		retry, err := GetRetry(ctx)
		if err != nil {
			return nil, err
		}
		producer, err := GetKafkaProducer(ctx)
		if err != nil {
			return nil, err
		}

		wrappedBatch := collect.NewCollect(
			collect.WithCollectCollector(c.Collector),
			collect.WithCollectAggregator(c.Aggregator),
			collect.WithCollectPersistenceIdFunc(c.PersistenceIdFunction),
		)

		wrappedBatch = stateless.NewProducerBatchFunction(
			stateless.WithBatchProducerNextFunction(wrappedBatch),
			stateless.WithBatchProducerRuntime(producer),
			stateless.WithBatchProducerPrometheus(),
		)

		wrappedBatch = stateless.NewBatchRetry(
			stateless.WithBatchRetryNextFunction(wrappedBatch),
			stateless.WithBatchRetryRuntime(retry),
			stateless.WithBatchRetryPrometheus(),
		)

		return wrappedBatch, nil
	})
}

func (c CollectorConfiguration) Runtime() runtime.Runtime {
	c.Register()

	return &RuntimeFacade{
		Runtimes: InjectedRuntimes(),
	}
}
