package flows

import (
	"context"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime_bun"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/inverse"
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

func (c StatefulPostgresqlFunctionConfiguration) Runtime() runtime.Runtime {
	RegisterPostgresql(c.PostgresqlConfiguration)
	RegisterPostgresqlSingleState(c.PersistenceTableName)
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
		repository, err := GetPostgresqlSingleStateRepository(ctx)
		if err != nil {
			return nil, err
		}

		wrappedStatefulFunction := c.StatefulFunction
		wrappedStatefulFunction = stateful.NewSingleStatefulDeduplicate(
			stateful.WithSingleStatefulDeduplicateNextFunction(wrappedStatefulFunction),
		)

		wrappedFunction := stateful.NewSingleReadWrite(
			stateful.WithSingleReadWriteStatefulFunction(wrappedStatefulFunction),
			stateful.WithSingleReadWriteTransactionPersistenceIdFunc(c.PersistenceIdFunction),
			stateful.WithSingleReadWriteRepository(repository),
		)

		wrappedFunction = stateless.NewSingleRetry(
			stateless.WithSingleRetryNextFunction(wrappedFunction),
			stateless.WithSingleRetryRuntime(retry),
			stateless.WithSingleRetryPrometheus(),
		)

		wrappedBatch := stateless.NewProducerBatchIterateFunction(
			stateless.WithBatchIterateFunctionNextFunction(wrappedFunction),
			stateless.WithBatchIterateFunctionProducer(producer),
			stateless.WithBatchIterateProducerPrometheus(),
		)

		wrappedBatch = stateless.NewBatchRetry(
			stateless.WithBatchRetryNextFunction(wrappedBatch),
			stateless.WithBatchRetryRuntime(retry),
			stateless.WithBatchRetryPrometheus(),
		)

		return wrappedBatch, nil
	})

	return &RuntimeFacade{
		Runtimes: InjectedRuntimes(),
	}
}
