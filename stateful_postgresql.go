package flows

import (
	"time"

	"github.com/avast/retry-go"
	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/stateful_bun"
	"github.com/hjwalt/flows/stateless"
)

// Wiring configuration
type StatefulPostgresqlFunctionConfiguration struct {
	PostgresqlConfiguration    []runtime.Configuration[*stateful_bun.PostgresqlConnection]
	KafkaProducerConfiguration []runtime.Configuration[*runtime_sarama.Producer]
	KafkaConsumerConfiguration []runtime.Configuration[*runtime_sarama.Consumer]
	StatefulFunction           stateful.SingleFunction
	PersistenceIdFunction      stateful.PersistenceIdFunction[[]byte, []byte]
	PersistenceTableName       string
	RouteConfiguration         []runtime.Configuration[*runtime_bunrouter.Router]
}

func (c StatefulPostgresqlFunctionConfiguration) Runtime() runtime.Runtime {

	ctrl := runtime.NewController()

	// postgres runtime
	postgresConnectionConfig := append(
		c.PostgresqlConfiguration,
		stateful_bun.WithController(ctrl),
	)
	conn := stateful_bun.NewPostgresqlConnection(postgresConnectionConfig...)

	// producer runtime
	producerConfig := append(
		c.KafkaProducerConfiguration,
		runtime_sarama.WithProducerRuntimeController(ctrl),
	)
	producer := runtime_sarama.NewProducer(producerConfig...)

	// bun transaction repository
	repository := stateful_bun.NewSingleStateRepository(
		stateful_bun.WithSingleStateRepositoryConnection(conn),
		stateful_bun.WithSingleStateRepositoryPersistenceTableName(c.PersistenceTableName),
	)

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
	messagesProduced := stateless.NewSingleProducer(
		stateless.WithSingleProducerNextFunction(stateTransaction),
		stateless.WithSingleProducerRuntime(producer),
		stateless.WithSingleProducerPrometheus(),
	)

	// - retry
	retryRuntime := runtime_retry.NewRetry(
		runtime_retry.WithRetryOption(
			retry.Attempts(1000000),
			retry.Delay(10*time.Millisecond),
			retry.MaxDelay(time.Second),
			retry.MaxJitter(time.Second),
			retry.DelayType(retry.BackOffDelay),
		),
	)
	produceRetry := stateless.NewSingleRetry(
		stateless.WithSingleRetryRuntime(retryRuntime),
		stateless.WithSingleRetryNextFunction(messagesProduced),
		stateless.WithSingleRetryPrometheus(),
	)

	// sarama consumer loop
	consumerLoop := runtime_sarama.NewSingleLoop(
		runtime_sarama.WithLoopSingleFunction(produceRetry),
		runtime_sarama.WithLoopSinglePrometheus(),
	)

	// consumer runtime
	consumerConfig := append(
		c.KafkaConsumerConfiguration,
		runtime_sarama.WithConsumerRuntimeController(ctrl),
		runtime_sarama.WithConsumerLoop(consumerLoop),
	)
	consumer := runtime_sarama.NewConsumer(consumerConfig...)

	// http runtime, prometheus first for hard prometheus path
	routeConfig := append(
		make([]runtime.Configuration[*runtime_bunrouter.Router], 0),
		runtime_bunrouter.WithRouterPrometheus(),
	)
	routeConfig = append(
		routeConfig,
		c.RouteConfiguration...,
	)
	routerRuntime := runtime_bunrouter.NewRouter(routeConfig...)

	// multi runtime configuration
	multi := runtime.NewMulti(
		runtime.WithController(ctrl),
		runtime.WithRuntime(routerRuntime),
		runtime.WithRuntime(conn),
		runtime.WithRuntime(producer),
		runtime.WithRuntime(consumer),
		runtime.WithRuntime(retryRuntime),
	)
	return multi
}
