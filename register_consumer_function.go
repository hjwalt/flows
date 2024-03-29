package flows

import (
	"context"

	"github.com/hjwalt/flows/collect"
	"github.com/hjwalt/flows/join"
	"github.com/hjwalt/flows/materialise"
	"github.com/hjwalt/flows/materialise/materialise_bun"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/stateful/stateful_batch_state"
	"github.com/hjwalt/flows/stateful/stateful_bun"
	"github.com/hjwalt/flows/stateful/stateful_deduplicate_offset"
	"github.com/hjwalt/flows/stateful/stateful_error_handler"
	"github.com/hjwalt/flows/stateful/stateful_error_handler_skip_list"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/flows/stateless/stateless_error_handler"
	"github.com/hjwalt/flows/stateless/stateless_error_handler_skip_list"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/structure"
)

func RegisterStatelessSingleFunction(
	ci inverse.Container,
	topic string,
	fn stateless.SingleFunction,
) {
	RegisterKafkaConsumerFunctionInjector(
		func(ctx context.Context, ci inverse.Container) (KafkaConsumerFunction, error) {
			retry, err := GetRetry(ctx, ci)
			if err != nil {
				return KafkaConsumerFunction{}, err
			}

			wrappedFunction := fn

			wrappedFunction = stateless_error_handler.New(
				wrappedFunction,
				stateless_error_handler_skip_list.Default(),
			)

			wrappedFunction = stateless.NewSingleRetry(
				stateless.WithSingleRetryNextFunction(wrappedFunction),
				stateless.WithSingleRetryRuntime(retry),
				stateless.WithSingleRetryPrometheus(),
			)

			wrappedBatch := stateless.NewBatchIterateFunction(
				stateless.WithBatchIterateFunctionNextFunction(wrappedFunction),
			)

			return KafkaConsumerFunction{
				Topic: topic,
				Key:   stateless.Base64PersistenceId,
				Fn:    wrappedBatch,
			}, nil
		},
		ci,
	)
}

func RegisterStatelessBatchFunction(
	ci inverse.Container,
	topic string,
	fn stateless.BatchFunction,
) {
	RegisterKafkaConsumerFunctionInjector(
		func(ctx context.Context, ci inverse.Container) (KafkaConsumerFunction, error) {
			return KafkaConsumerFunction{
				Topic: topic,
				Key:   stateless.Base64PersistenceId,
				Fn:    fn,
			}, nil
		},
		ci,
	)
}

func RegisterStatelessSingleFunctionWithKey(
	ci inverse.Container,
	topic string,
	fn stateless.SingleFunction,
	key stateful.PersistenceIdFunction[structure.Bytes, structure.Bytes],
) {
	RegisterKafkaConsumerFunctionInjector(
		func(ctx context.Context, ci inverse.Container) (KafkaConsumerFunction, error) {
			retry, err := GetRetry(ctx, ci)
			if err != nil {
				return KafkaConsumerFunction{}, err
			}

			wrappedFunction := fn

			wrappedFunction = stateless_error_handler.New(
				wrappedFunction,
				stateless_error_handler_skip_list.Default(),
			)

			wrappedFunction = stateless.NewSingleRetry(
				stateless.WithSingleRetryNextFunction(wrappedFunction),
				stateless.WithSingleRetryRuntime(retry),
				stateless.WithSingleRetryPrometheus(),
			)

			wrappedBatch := stateless.NewBatchIterateFunction(
				stateless.WithBatchIterateFunctionNextFunction(wrappedFunction),
			)

			return KafkaConsumerFunction{
				Topic: topic,
				Key:   key,
				Fn:    wrappedBatch,
			}, nil
		},
		ci,
	)
}

func RegisterStatefulFunction(
	ci inverse.Container,
	topic string,
	tableName string,
	fn stateful.SingleFunction,
	key stateful.PersistenceIdFunction[structure.Bytes, structure.Bytes],
) {
	RegisterKafkaConsumerFunctionInjector(
		func(ctx context.Context, ci inverse.Container) (KafkaConsumerFunction, error) {
			bunConnection, getBunConnectionError := GetPostgresqlConnection(ctx, ci)
			if getBunConnectionError != nil {
				return KafkaConsumerFunction{}, getBunConnectionError
			}

			repository := stateful_bun.NewRepository(
				stateful_bun.WithConnection(bunConnection),
				stateful_bun.WithStateTableName(tableName),
			)

			wrappedStatefulFunction := fn

			wrappedStatefulFunction = stateful_error_handler.New(
				wrappedStatefulFunction,
				stateful_error_handler_skip_list.Default(),
			)

			wrappedStatefulFunction = stateful_deduplicate_offset.New(
				wrappedStatefulFunction,
			)

			wrappedBatch := stateful_batch_state.New(
				wrappedStatefulFunction,
				key,
				repository,
			)

			return KafkaConsumerFunction{
				Topic: topic,
				Key:   key,
				Fn:    wrappedBatch,
			}, nil
		},
		ci,
	)
}

func RegisterJoinStatefulFunction(
	ci inverse.Container,
	topic string,
	tableName string,
	fn stateful.SingleFunction,
	key stateful.PersistenceIdFunction[structure.Bytes, structure.Bytes],
) {
	RegisterKafkaConsumerFunctionInjector(
		func(ctx context.Context, ci inverse.Container) (KafkaConsumerFunction, error) {
			bunConnection, getBunConnectionError := GetPostgresqlConnection(ctx, ci)
			if getBunConnectionError != nil {
				return KafkaConsumerFunction{}, getBunConnectionError
			}

			repository := stateful_bun.NewRepository(
				stateful_bun.WithConnection(bunConnection),
				stateful_bun.WithStateTableName(tableName),
			)

			wrappedStatefulFunction := fn

			wrappedStatefulFunction = stateful_deduplicate_offset.New(
				wrappedStatefulFunction,
			)

			wrappedBatch := stateful_batch_state.New(
				wrappedStatefulFunction,
				key,
				repository,
			)

			wrappedBatch = join.NewIntermediateToJoinMap(
				join.WithIntermediateToJoinMapTransactionWrappedFunction(wrappedBatch),
			)

			return KafkaConsumerFunction{
				Topic: topic,
				Key:   join.IntermediateTopicKeyFunction,
				Fn:    wrappedBatch,
			}, nil
		},
		ci,
	)
}

func RegisterMaterialiseFunction[T any](
	ci inverse.Container,
	topic string,
	fn materialise.MapFunction[structure.Bytes, structure.Bytes, T],
) {
	RegisterKafkaConsumerFunctionInjector(
		func(ctx context.Context, ci inverse.Container) (KafkaConsumerFunction, error) {
			bunConnection, getBunConnectionError := GetPostgresqlConnection(ctx, ci)
			if getBunConnectionError != nil {
				return KafkaConsumerFunction{}, getBunConnectionError
			}

			repository := materialise_bun.NewBunUpsertRepository(
				materialise_bun.WithBunUpsertRepositoryConnection[T](bunConnection),
			)

			wrappedBatch := materialise.NewBatchUpsert(
				materialise.WithBatchUpsertMapFunction(fn),
				materialise.WithBatchUpsertRepository[T](repository),
			)

			return KafkaConsumerFunction{
				Topic: topic,
				Key:   stateless.Base64PersistenceId,
				Fn:    wrappedBatch,
			}, nil
		},
		ci,
	)
}

func RegisterCollectorFunction(
	ci inverse.Container,
	topic string,
	persistenceIdFunction stateful.PersistenceIdFunction[[]byte, []byte],
	aggregator collect.Aggregator[structure.Bytes, structure.Bytes, structure.Bytes],
	collector collect.Collector,
) {
	RegisterKafkaConsumerFunctionInjector(
		func(ctx context.Context, ci inverse.Container) (KafkaConsumerFunction, error) {
			wrappedBatch := collect.NewCollect(
				collect.WithCollectCollector(collector),
				collect.WithCollectAggregator(aggregator),
				collect.WithCollectPersistenceIdFunc(persistenceIdFunction),
			)

			return KafkaConsumerFunction{
				Topic: topic,
				Key:   stateless.Base64PersistenceId,
				Fn:    wrappedBatch,
			}, nil
		},
		ci,
	)
}
