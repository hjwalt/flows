package flows

import (
	"github.com/hjwalt/flows/materialise"
	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime_bun"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/runway/runtime"
)

type MaterialisePostgresqlFunctionConfiguration[T any] struct {
	PostgresqlConfiguration    []runtime.Configuration[*runtime_bun.PostgresqlConnection]
	KafkaProducerConfiguration []runtime.Configuration[*runtime_sarama.Producer]
	KafkaConsumerConfiguration []runtime.Configuration[*runtime_sarama.Consumer]
	RetryConfiguration         []runtime.Configuration[*runtime_retry.Retry]
	MaterialiseMapFunction     materialise.MapFunction[message.Bytes, message.Bytes, T]
	RouteConfiguration         []runtime.Configuration[*runtime_bunrouter.Router]
	AdditionalRuntimes         []runtime.Runtime
}

func (c MaterialisePostgresqlFunctionConfiguration[T]) Runtime() runtime.Runtime {

	// postgres runtime
	conn := Postgresql(c.PostgresqlConfiguration)
	repository := PostgresqlUpsertRepository[T](conn)

	// producer runtime
	producer := KafkaProducer(c.KafkaProducerConfiguration)

	// function
	materialiseFn := materialise.NewSingleUpsert(
		materialise.WithSingleUpsertRepository(repository),
		materialise.WithSingleUpsertMapFunction(c.MaterialiseMapFunction),
	)

	retriedFn, retryRuntime := WrapRetry(materialiseFn, c.RetryConfiguration)

	// consumer runtime
	consumer := KafkaConsumerSingle(retriedFn, c.KafkaConsumerConfiguration)

	// http runtime
	routerRuntime := RouteRuntime(producer, c.RouteConfiguration)

	// add additional runtimes
	runtimes := []runtime.Runtime{
		conn,
		routerRuntime,
		producer,
		consumer,
		retryRuntime,
	}
	if len(c.AdditionalRuntimes) > 0 {
		runtimes = append(c.AdditionalRuntimes, runtimes...)
	}

	return &RuntimeFacade{
		Runtimes: runtimes,
	}
}
