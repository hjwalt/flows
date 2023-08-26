package flows

import (
	"context"

	"github.com/hjwalt/flows/materialise"
	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime_bun"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateless"
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

func (c MaterialisePostgresqlFunctionConfiguration[T]) Register() {
	RegisterPostgresqlConfig(c.PostgresqlConfiguration...)
	RegisterPostgresql()
	RegisterPostgresqlUpsert[T]()
	RegisterRetry(c.RetryConfiguration)
	RegisterProducerConfig(c.KafkaProducerConfiguration...)
	RegisterProducer()
	RegisterConsumerConfig(c.KafkaConsumerConfiguration...)
	RegisterConsumerKeyedConfig()
	RegisterConsumer()
	RegisterRouteConfigDefault()
	RegisterRouteConfig(c.RouteConfiguration...)
	RegisterRoute()
	RegisterConsumerKeyedKeyFunction(stateless.Base64PersistenceId)
	RegisterConsumerKeyedFunction(func(ctx context.Context) (stateless.BatchFunction, error) {
		retry, err := GetRetry(ctx)
		if err != nil {
			return nil, err
		}
		repository, err := GetPostgresqlUpsertRepository[T](ctx)
		if err != nil {
			return nil, err
		}

		wrappedBatch := materialise.NewBatchUpsert(
			materialise.WithBatchUpsertMapFunction(c.MaterialiseMapFunction),
			materialise.WithBatchUpsertRepository(repository),
		)

		wrappedBatch = stateless.NewBatchRetry(
			stateless.WithBatchRetryNextFunction(wrappedBatch),
			stateless.WithBatchRetryRuntime(retry),
			stateless.WithBatchRetryPrometheus(),
		)

		return wrappedBatch, nil
	})
}

func (c MaterialisePostgresqlFunctionConfiguration[T]) Runtime() runtime.Runtime {
	c.Register()

	return &RuntimeFacade{
		Runtimes: InjectedRuntimes(),
	}
}
