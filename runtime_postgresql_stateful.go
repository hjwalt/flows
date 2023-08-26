package flows

import (
	"context"

	"github.com/hjwalt/flows/runtime_bun"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/runtime"
)

// Wiring configuration
type StatefulPostgresqlFunctionConfiguration struct {
	PostgresqlConfiguration    []runtime.Configuration[*runtime_bun.PostgresqlConnection]
	KafkaProducerConfiguration []runtime.Configuration[*runtime_sarama.Producer]
	KafkaConsumerConfiguration []runtime.Configuration[*runtime_sarama.Consumer]
	RetryConfiguration         []runtime.Configuration[*runtime_retry.Retry]
	StatefulFunction           stateful.SingleFunction
	PersistenceIdFunction      stateful.PersistenceIdFunction[[]byte, []byte]
	PersistenceTableName       string
	RouteConfiguration         []runtime.Configuration[*runtime_bunrouter.Router]
}

func (c StatefulPostgresqlFunctionConfiguration) Register() {
	RegisterPostgresqlConfig(c.PostgresqlConfiguration...)
	RegisterPostgresql()
	RegisterPostgresqlSingleState(c.PersistenceTableName)
	RegisterRetry(c.RetryConfiguration)
	RegisterProducerConfig(c.KafkaProducerConfiguration...)
	RegisterProducer()
	RegisterConsumerConfig(c.KafkaConsumerConfiguration...)
	RegisterConsumerKeyedConfig()
	RegisterConsumer()
	RegisterRouteConfigDefault()
	RegisterRouteConfig(c.RouteConfiguration...)
	RegisterRoute()
	RegisterConsumerKeyedKeyFunction(c.PersistenceIdFunction)
	RegisterConsumerKeyedFunction(func(ctx context.Context) (stateless.BatchFunction, error) {
		retry, err := GetRetry(ctx)
		if err != nil {
			return nil, err
		}
		producer, err := GetKafkaProducer(ctx)
		if err != nil {
			return nil, err
		}
		repository, err := GetPostgresqlSingleStateRepository(ctx)
		if err != nil {
			return nil, err
		}

		wrappedStatefulFunction := c.StatefulFunction
		wrappedStatefulFunction = stateful.NewDeduplicate(
			stateful.WithDeduplicateNextFunction(wrappedStatefulFunction),
		)

		wrappedBatch := stateful.NewReadWrite(
			stateful.WithReadWriteFunction(wrappedStatefulFunction),
			stateful.WithReadWritePersistenceIdFunc(c.PersistenceIdFunction),
			stateful.WithReadWriteRepository(repository),
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

func (c StatefulPostgresqlFunctionConfiguration) Runtime() runtime.Runtime {
	c.Register()

	return &RuntimeFacade{
		Runtimes: InjectedRuntimes(),
	}
}
