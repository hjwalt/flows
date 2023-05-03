package flows

import (
	"github.com/hjwalt/flows/materialise"
	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/flows/runtime_bun"
	"github.com/hjwalt/flows/runtime_sarama"
)

type MaterialisePostgresqlFunctionConfiguration[T any] struct {
	PostgresqlConfiguration    []runtime.Configuration[*runtime_bun.PostgresqlConnection]
	KafkaConsumerConfiguration []runtime.Configuration[*runtime_sarama.Consumer]
	MaterialiseMapFunction     materialise.MapFunction[T]
}

func (c MaterialisePostgresqlFunctionConfiguration[T]) Runtime() runtime.Runtime {
	ctrl := runtime.NewController()

	// postgres runtime
	conn := Postgresql(ctrl, c.PostgresqlConfiguration)
	repository := PostgresqlUpsertRepository[T](conn)

	// function
	materialiseFn := materialise.NewSingleUpsert[T](
		materialise.WithSingleUpsertRepository[T](repository),
		materialise.WithSingleUpsertMapFunction[T](c.MaterialiseMapFunction),
	)

	retriedFn, retryRuntime := WrapRetry(materialiseFn)

	consumer := KafkaConsumerSingle(ctrl, retriedFn, c.KafkaConsumerConfiguration)

	// multi runtime configuration
	multi := runtime.NewMulti(
		runtime.WithController(ctrl),
		runtime.WithRuntime(conn),
		runtime.WithRuntime(consumer),
		runtime.WithRuntime(retryRuntime),
	)
	return multi
}
