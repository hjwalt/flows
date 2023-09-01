package flows

import (
	"context"

	"github.com/hjwalt/flows/collect"
	"github.com/hjwalt/flows/materialise"
	"github.com/hjwalt/flows/materialise_bun"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/stateful_bun"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/structure"
)

func RegisterStatelessSingleFunction(topic string, fn stateless.SingleFunction) {
	RegisterConsumerFunctionInjector(
		func(ctx context.Context) (ConsumerFunction, error) {
			retry, err := GetRetry(ctx)
			if err != nil {
				return ConsumerFunction{}, err
			}

			wrappedFunction := fn

			wrappedFunction = stateless.NewSingleRetry(
				stateless.WithSingleRetryNextFunction(wrappedFunction),
				stateless.WithSingleRetryRuntime(retry),
				stateless.WithSingleRetryPrometheus(),
			)

			wrappedBatch := stateless.NewBatchIterateFunction(
				stateless.WithBatchIterateFunctionNextFunction(wrappedFunction),
			)

			return ConsumerFunction{
				Topic: topic,
				Key:   stateless.Base64PersistenceId,
				Fn:    wrappedBatch,
			}, nil
		},
	)
}

func RegisterStatelessSingleFunctionWithKey(topic string, fn stateless.SingleFunction, key stateful.PersistenceIdFunction[structure.Bytes, structure.Bytes]) {
	RegisterConsumerFunctionInjector(
		func(ctx context.Context) (ConsumerFunction, error) {
			retry, err := GetRetry(ctx)
			if err != nil {
				return ConsumerFunction{}, err
			}

			wrappedFunction := fn

			wrappedFunction = stateless.NewSingleRetry(
				stateless.WithSingleRetryNextFunction(wrappedFunction),
				stateless.WithSingleRetryRuntime(retry),
				stateless.WithSingleRetryPrometheus(),
			)

			wrappedBatch := stateless.NewBatchIterateFunction(
				stateless.WithBatchIterateFunctionNextFunction(wrappedFunction),
			)

			return ConsumerFunction{
				Topic: topic,
				Key:   key,
				Fn:    wrappedBatch,
			}, nil
		},
	)
}

func RegisterStatefulFunction(
	topic string,
	tableName string,
	fn stateful.SingleFunction,
	key stateful.PersistenceIdFunction[structure.Bytes, structure.Bytes],
) {
	RegisterConsumerFunctionInjector(
		func(ctx context.Context) (ConsumerFunction, error) {
			bunConnection, getBunConnectionError := GetPostgresqlConnection(ctx)
			if getBunConnectionError != nil {
				return ConsumerFunction{}, getBunConnectionError
			}

			repository := stateful_bun.NewRepository(
				stateful_bun.WithConnection(bunConnection),
				stateful_bun.WithStateTableName(tableName),
			)

			wrappedStatefulFunction := fn
			wrappedStatefulFunction = stateful.NewDeduplicate(
				stateful.WithDeduplicateNextFunction(wrappedStatefulFunction),
			)

			wrappedBatch := stateful.NewReadWrite(
				stateful.WithReadWriteFunction(wrappedStatefulFunction),
				stateful.WithReadWritePersistenceIdFunc(key),
				stateful.WithReadWriteRepository(repository),
			)

			return ConsumerFunction{
				Topic: topic,
				Key:   key,
				Fn:    wrappedBatch,
			}, nil

		},
	)
}

func RegisterMaterialiseFunction[T any](topic string, fn materialise.MapFunction[structure.Bytes, structure.Bytes, T]) {
	RegisterConsumerFunctionInjector(
		func(ctx context.Context) (ConsumerFunction, error) {
			bunConnection, getBunConnectionError := GetPostgresqlConnection(ctx)
			if getBunConnectionError != nil {
				return ConsumerFunction{}, getBunConnectionError
			}

			repository := materialise_bun.NewBunUpsertRepository(
				materialise_bun.WithBunUpsertRepositoryConnection[T](bunConnection),
			)

			wrappedBatch := materialise.NewBatchUpsert(
				materialise.WithBatchUpsertMapFunction(fn),
				materialise.WithBatchUpsertRepository[T](repository),
			)

			return ConsumerFunction{
				Topic: topic,
				Key:   stateless.Base64PersistenceId,
				Fn:    wrappedBatch,
			}, nil
		},
	)
}

func RegisterCollectorFunction(
	topic string,
	persistenceIdFunction stateful.PersistenceIdFunction[[]byte, []byte],
	aggregator collect.Aggregator[structure.Bytes, structure.Bytes, structure.Bytes],
	collector collect.Collector,
) {
	RegisterConsumerFunctionInjector(
		func(ctx context.Context) (ConsumerFunction, error) {
			wrappedBatch := collect.NewCollect(
				collect.WithCollectCollector(collector),
				collect.WithCollectAggregator(aggregator),
				collect.WithCollectPersistenceIdFunc(persistenceIdFunction),
			)

			return ConsumerFunction{
				Topic: topic,
				Key:   stateless.Base64PersistenceId,
				Fn:    wrappedBatch,
			}, nil
		},
	)
}
