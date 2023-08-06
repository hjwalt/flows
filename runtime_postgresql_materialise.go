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
	RegisterConsumerSingleConfig(c.KafkaConsumerConfiguration)
	RegisterConsumerSingle()
	RegisterRoute(c.RouteConfiguration)
	inverse.Register[stateless.SingleFunction](QualifierKafkaConsumerSingleFunction, func(ctx context.Context) (stateless.SingleFunction, error) {
		retry, err := GetRetry(ctx)
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

		return wrappedFunction, nil
	})

	return &RuntimeFacade{
		Runtimes: InjectedRuntimes(),
	}
}
