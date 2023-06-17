package flows

import (
	"github.com/hjwalt/flows/materialise"
	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/flows/runtime_bun"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_sarama"
)

type MaterialisePostgresqlFunctionConfiguration[T any] struct {
	PostgresqlConfiguration    []runtime.Configuration[*runtime_bun.PostgresqlConnection]
	KafkaProducerConfiguration []runtime.Configuration[*runtime_sarama.Producer]
	KafkaConsumerConfiguration []runtime.Configuration[*runtime_sarama.Consumer]
	MaterialiseMapFunction     materialise.MapFunction[message.Bytes, message.Bytes, T]
	RouteConfiguration         []runtime.Configuration[*runtime_bunrouter.Router]
}

func (c MaterialisePostgresqlFunctionConfiguration[T]) Runtime() runtime.Runtime {
	ctrl := runtime.NewController()

	// postgres runtime
	conn := Postgresql(ctrl, c.PostgresqlConfiguration)
	repository := PostgresqlUpsertRepository[T](conn)

	// producer runtime
	producer := KafkaProducer(ctrl, c.KafkaProducerConfiguration)

	// function
	materialiseFn := materialise.NewSingleUpsert(
		materialise.WithSingleUpsertRepository(repository),
		materialise.WithSingleUpsertMapFunction(c.MaterialiseMapFunction),
	)

	retriedFn, retryRuntime := WrapRetry(materialiseFn)

	// consumer runtime
	consumer := KafkaConsumerSingle(ctrl, retriedFn, c.KafkaConsumerConfiguration)

	// http runtime
	routerRuntime := RouteRuntime(producer, c.RouteConfiguration)

	// multi runtime configuration
	multi := runtime.NewMulti(
		runtime.WithController(ctrl),
		runtime.WithRuntime(conn),
		runtime.WithRuntime(routerRuntime),
		runtime.WithRuntime(producer),
		runtime.WithRuntime(consumer),
		runtime.WithRuntime(retryRuntime),
	)
	return multi
}
