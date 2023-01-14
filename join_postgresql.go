package flows

import (
	"time"

	"github.com/avast/retry-go"
	"github.com/hjwalt/flows/join"
	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/stateful_bun"
	"github.com/hjwalt/flows/stateless"
)

// Wiring configuration
// Hard to describe without a execution tree graph, but the rough idea is as follows
// consumer -> retry
// retry -> producer
// producer -> stateless switch
// stateless switch -> source to intermediate (for each source topic)
// stateless switch -> intermediate to join
// intermediate to join -> transaction
// transaction -> offset deduplication
// offset deduplication -> stateful switch
// stateful switch -> stateful function (for each source topic)

type JoinPostgresqlFunctionConfiguration struct {
	PostgresqlConfiguration    []runtime.Configuration[*stateful_bun.PostgresqlConnection]
	KafkaProducerConfiguration []runtime.Configuration[*runtime_sarama.Producer]
	KafkaConsumerConfiguration []runtime.Configuration[*runtime_sarama.Consumer]
	StatefulFunctions          map[string]stateful.StatefulBinarySingleFunction
	PersistenceIdFunctions     map[string]stateful.StatefulBinaryPersistenceIdFunction
	IntermediateTopicName      string
	PersistenceTableName       string
	RouteConfiguration         []runtime.Configuration[*runtime_bunrouter.Router]
}

func (c JoinPostgresqlFunctionConfiguration) Runtime() runtime.Runtime {

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

	topics := []string{}
	statefulTopicSwitchConfigurations := []runtime.Configuration[*stateful.SingleTopicSwitch]{}
	persistenceIdConfigurations := []runtime.Configuration[*stateful.PersistenceIdSwitch]{}

	// generatic stateful transaction
	for topic, statefulFn := range c.StatefulFunctions {
		persistenceIdFn, persistenceIdFnExists := c.PersistenceIdFunctions[topic]
		if !persistenceIdFnExists {
			// TODO: not so silent failure
			continue
		}

		topics = append(topics, topic)
		statefulTopicSwitchConfigurations = append(statefulTopicSwitchConfigurations, stateful.WithSingleTopicSwitchStatefulSingleFunction(topic, statefulFn))
		persistenceIdConfigurations = append(persistenceIdConfigurations, stateful.WithPersistenceIdSwitchPersistenceIdFunction(topic, persistenceIdFn))
	}

	// stateful topic switch and persistence id switch
	statefulTopicSwitch := stateful.NewSingleTopicSwitch(statefulTopicSwitchConfigurations...)
	persistenceIdTopicSwitch := stateful.NewSinglePersistenceIdSwitch(persistenceIdConfigurations...)

	// stateful function switch
	// - offset deduplication
	offsetDeduplicated := stateful.NewSingleStatefulDeduplicate(
		stateful.WithSingleStatefulDeduplicateNextFunction(statefulTopicSwitch),
	)

	// - transaction with bun
	stateTransaction := stateful.NewSingleReadWrite(
		stateful.WithSingleReadWriteTransactionPersistenceIdFunc(persistenceIdTopicSwitch),
		stateful.WithSingleReadWriteRepository(repository),
		stateful.WithSingleReadWriteStatefulFunction(offsetDeduplicated),
	)

	// - intermediate to join function
	intermediateToJoin := join.NewIntermediateToJoinMap(
		join.WithIntermediateToJoinMapTransactionWrappedFunction(stateTransaction),
	)

	// generating source to intermediate join
	sourceToIntermediateMap := join.NewSourceToIntermediateMap(
		join.WithSourceToIntermediateMapIntermediateTopic(c.IntermediateTopicName),
		join.WithSourceToIntermediateMapPersistenceIdFunction(persistenceIdTopicSwitch),
	)

	// generating stateless topic switch
	statelessTopicSwitchConfigurations := []runtime.Configuration[*stateless.SingleTopicSwitch]{
		stateless.WithSingleTopicSwitchStatelessSingleFunction(c.IntermediateTopicName, intermediateToJoin),
	}

	for _, topic := range topics {
		statelessTopicSwitchConfigurations = append(statelessTopicSwitchConfigurations, stateless.WithSingleTopicSwitchStatelessSingleFunction(topic, sourceToIntermediateMap))
	}

	statelessTopicSwitch := stateless.NewSingleTopicSwitch(
		statelessTopicSwitchConfigurations...,
	)

	// stateless topic switch wrapping
	// - produce output messages
	messagesProduced := stateless.NewSingleProducer(
		stateless.WithSingleProducerNextFunction(statelessTopicSwitch),
		stateless.WithSingleProducerRuntime(producer),
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
	)

	// sarama consumer loop
	consumerLoop := runtime_sarama.NewSingleLoop(
		runtime_sarama.WithLoopSingleFunction(produceRetry),
	)

	// consumer runtime
	topics = append(topics, c.IntermediateTopicName)

	consumerConfig := append(
		c.KafkaConsumerConfiguration,
		runtime_sarama.WithConsumerRuntimeController(ctrl),
		runtime_sarama.WithConsumerLoop(consumerLoop),
		runtime_sarama.WithConsumerTopic(topics...),
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
