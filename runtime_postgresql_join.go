package flows

import (
	"context"

	"github.com/hjwalt/flows/join"
	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime_bun"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/inverse"
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

	topics := []string{}
	statefulTopicSwitchConfigurations := []runtime.Configuration[*stateful.SingleTopicSwitch]{}
	persistenceIdConfigurations := []runtime.Configuration[*stateful.PersistenceIdSwitch]{}

	// generating stateful switches for persistence id and stateful function
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
	keyedPersistenceIdConfigurations := append(persistenceIdConfigurations, stateful.WithPersistenceIdSwitchPersistenceIdFunction(c.IntermediateTopicName, intermediateTopicKeyFunction))
	keyedPersistenceId := stateful.NewSinglePersistenceIdSwitch(keyedPersistenceIdConfigurations...)

	// setting the topics for consumers
	consumerTopics := append(topics, c.IntermediateTopicName)
	inverse.RegisterInstance[runtime.Configuration[*runtime_sarama.Consumer]](QualifierKafkaConsumerConfiguration, runtime_sarama.WithConsumerTopic(consumerTopics...))

	RegisterPostgresql(c.PostgresqlConfiguration)
	RegisterPostgresqlSingleState(c.PersistenceTableName)
	RegisterRetry(c.RetryConfiguration)
	RegisterProducerConfig(c.KafkaProducerConfiguration)
	RegisterProducer()
	RegisterConsumerKeyedConfig(c.KafkaConsumerConfiguration)
	RegisterConsumer()
	RegisterRoute(c.RouteConfiguration)
	inverse.RegisterInstance[stateful.PersistenceIdFunction[message.Bytes, message.Bytes]](QualifierKafkaConsumerKeyFunction, keyedPersistenceId)
	inverse.Register[stateless.BatchFunction](QualifierKafkaConsumerBatchFunction, func(ctx context.Context) (stateless.BatchFunction, error) {
		retry, err := GetRetry(ctx)
		if err != nil {
			return nil, err
		}
		producer, err := GetKafkaProducer(ctx)
		if err != nil {
			return nil, err
		}
		repository, err := GetPostgresqlSingleStateRepository(ctx)
		if err != nil {
			return nil, err
		}

		// Intermediate to join stateful switching function
		statefulWrappedFunction := stateful.NewSingleTopicSwitch(statefulTopicSwitchConfigurations...)
		persistenceIdTopicSwitch := stateful.NewSinglePersistenceIdSwitch(persistenceIdConfigurations...)

		statefulWrappedFunction = stateful.NewSingleStatefulDeduplicate(
			stateful.WithSingleStatefulDeduplicateNextFunction(statefulWrappedFunction),
		)

		stateTransaction := stateful.NewBatchReadWrite(
			stateful.WithBatchReadWritePersistenceIdFunc(persistenceIdTopicSwitch),
			stateful.WithBatchReadWriteRepository(repository),
			stateful.WithBatchReadWriteStatefulFunction(statefulWrappedFunction),
		)

		intermediateToJoin := join.NewIntermediateToJoinMap(
			join.WithIntermediateToJoinMapTransactionWrappedFunction(stateTransaction),
		)

		// Stateless topic switch

		// generating source to intermediate join to be added into stateless topic switch
		sourceToIntermediateMap := join.NewSourceToIntermediateMap(
			join.WithSourceToIntermediateMapIntermediateTopic(c.IntermediateTopicName),
			join.WithSourceToIntermediateMapPersistenceIdFunction(persistenceIdTopicSwitch),
		)

		topicSwitch := join.NewJoinSwitch(
			join.WithJoinSwitchIntermediateTopicFunction(c.IntermediateTopicName, intermediateToJoin),
			join.WithJoinSwitchSourceTopicFunction(sourceToIntermediateMap),
		)

		wrappedBatch := stateless.NewProducerBatchFunction(
			stateless.WithBatchProducerNextFunction(topicSwitch),
			stateless.WithBatchProducerRuntime(producer),
			stateless.WithBatchProducerPrometheus(),
		)

		wrappedBatch = stateless.NewBatchRetry(
			stateless.WithBatchRetryNextFunction(wrappedBatch),
			stateless.WithBatchRetryRuntime(retry),
			stateless.WithBatchRetryPrometheus(),
		)

		return wrappedBatch, nil
	})

	return &RuntimeFacade{
		Runtimes: InjectedRuntimes(),
	}
}

func intermediateTopicKeyFunction(ctx context.Context, m message.Message[message.Bytes, message.Bytes]) (string, error) {
	keyValue, keyError := join.IntermediateKeyFormat.Unmarshal(m.Key)
	if keyError != nil {
		return "", keyError
	}
	return keyValue.PersistenceId, nil
}
