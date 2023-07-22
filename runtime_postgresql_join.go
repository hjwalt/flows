package flows

import (
	"github.com/hjwalt/flows/join"
	"github.com/hjwalt/flows/runtime_bun"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/runtime"
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
// stateful switch -> stateful function(s) (for each source topic)

type JoinPostgresqlFunctionConfiguration struct {
	PostgresqlConfiguration    []runtime.Configuration[*runtime_bun.PostgresqlConnection]
	KafkaProducerConfiguration []runtime.Configuration[*runtime_sarama.Producer]
	KafkaConsumerConfiguration []runtime.Configuration[*runtime_sarama.Consumer]
	RetryConfiguration         []runtime.Configuration[*runtime_retry.Retry]
	StatefulFunctions          map[string]stateful.SingleFunction
	PersistenceIdFunctions     map[string]stateful.PersistenceIdFunction[[]byte, []byte]
	IntermediateTopicName      string
	PersistenceTableName       string
	RouteConfiguration         []runtime.Configuration[*runtime_bunrouter.Router]
}

func (c JoinPostgresqlFunctionConfiguration) Runtime() runtime.Runtime {

	// postgres runtime
	conn := Postgresql(c.PostgresqlConfiguration)
	repository := PostgresqlSingleStateRepository(conn, c.PersistenceTableName)

	// producer runtime
	producer := KafkaProducer(c.KafkaProducerConfiguration)

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
	messagesProduced := WrapSingleProduce(statelessTopicSwitch, producer)

	// - retry
	produceRetry, retryRuntime := WrapRetry(messagesProduced, c.RetryConfiguration)

	// sarama consumer loop
	consumerLoop := runtime_sarama.NewSingleLoop(
		runtime_sarama.WithLoopSingleFunction(produceRetry),
		runtime_sarama.WithLoopSinglePrometheus(),
	)

	// consumer runtime
	topics = append(topics, c.IntermediateTopicName)

	consumerConfig := append(
		c.KafkaConsumerConfiguration,
		runtime_sarama.WithConsumerLoop(consumerLoop),
		runtime_sarama.WithConsumerTopic(topics...),
	)
	consumer := runtime_sarama.NewConsumer(consumerConfig...)

	// http runtime
	routerRuntime := RouteRuntime(producer, c.RouteConfiguration)

	return &RuntimeFacade{
		Runtimes: []runtime.Runtime{
			conn,
			routerRuntime,
			producer,
			consumer,
			retryRuntime,
		},
	}
}
