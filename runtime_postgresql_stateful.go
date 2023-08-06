package flows

import (
	"github.com/hjwalt/flows/runtime_bun"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateful"
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
	AdditionalRuntimes         []runtime.Runtime
}

func (c StatefulPostgresqlFunctionConfiguration) Runtime() runtime.Runtime {

	// postgres runtime
	conn := Postgresql(c.PostgresqlConfiguration)
	repository := PostgresqlSingleStateRepository(conn, c.PersistenceTableName)

	// producer runtime
	producer := KafkaProducer(c.KafkaProducerConfiguration)

	// function wrapping
	// - offset deduplication
	offsetDeduplicated := stateful.NewSingleStatefulDeduplicate(
		stateful.WithSingleStatefulDeduplicateNextFunction(c.StatefulFunction),
	)

	// - transaction with bun
	stateTransaction := stateful.NewSingleReadWrite(
		stateful.WithSingleReadWriteTransactionPersistenceIdFunc(c.PersistenceIdFunction),
		stateful.WithSingleReadWriteRepository(repository),
		stateful.WithSingleReadWriteStatefulFunction(offsetDeduplicated),
	)

	// - produce output messages
	messagesProduced := WrapSingleProduce(stateTransaction, producer)

	// - retry
	produceRetry, retryRuntime := WrapRetry(messagesProduced, c.RetryConfiguration)

	// consumer runtime
	consumer := KafkaConsumerSingle(produceRetry, c.KafkaConsumerConfiguration)

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
