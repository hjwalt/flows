package flows

import (
	"context"

	"github.com/hjwalt/flows/materialise"
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

type MaterialisePostgresqlFunctionConfiguration[T any] struct {
	PostgresqlConfiguration    []runtime.Configuration[*runtime_bun.PostgresqlConnection]
	KafkaProducerConfiguration []runtime.Configuration[*runtime_sarama.Producer]
	KafkaConsumerConfiguration []runtime.Configuration[*runtime_sarama.Consumer]
	RetryConfiguration         []runtime.Configuration[*runtime_retry.Retry]
	MaterialiseMapFunction     materialise.MapFunction[message.Bytes, message.Bytes, T]
	RouteConfiguration         []runtime.Configuration[*runtime_bunrouter.Router]
}

func (c MaterialisePostgresqlFunctionConfiguration[T]) Runtime() runtime.Runtime {
	RegisterPostgresql(c.PostgresqlConfiguration)
	RegisterPostgresqlUpsert[T]()
	RegisterRetry(c.RetryConfiguration)
	RegisterProducerConfig(c.KafkaProducerConfiguration)
	RegisterProducer()
	RegisterConsumerKeyedConfig(c.KafkaConsumerConfiguration)
	RegisterConsumer()
	RegisterRoute(c.RouteConfiguration)
	inverse.RegisterInstance[stateful.PersistenceIdFunction[message.Bytes, message.Bytes]](QualifierKafkaConsumerKeyFunction, statelessBase64PersistenceId)
	inverse.Register[stateless.BatchFunction](QualifierKafkaConsumerBatchFunction, func(ctx context.Context) (stateless.BatchFunction, error) {
		retry, err := GetRetry(ctx)
		if err != nil {
			return nil, err
		}
		producer, err := GetKafkaProducer(ctx)
		if err != nil {
			return nil, err
		}
		repository, err := GetPostgresqlUpsertRepository[T](ctx)
		if err != nil {
			return nil, err
		}

		wrappedFunction := materialise.NewSingleUpsert(
			materialise.WithSingleUpsertMapFunction(c.MaterialiseMapFunction),
			materialise.WithSingleUpsertRepository(repository),
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
