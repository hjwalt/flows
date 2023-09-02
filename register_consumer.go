package flows

import (
	"context"
	"errors"

	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
	"github.com/hjwalt/runway/structure"
)

const (
	QualifierKafkaConsumerConfiguration             = "QualifierKafkaConsumerConfiguration"
	QualifierKafkaConsumer                          = "QualifierKafkaConsumer"
	QualifierKafkaConsumerHandler                   = "QualifierKafkaConsumerHandler"
	QualifierKafkaConsumerBatchFunction             = "QualifierKafkaConsumerBatchFunction"
	QualifierKafkaConsumerKeyFunction               = "QualifierKafkaConsumerKeyFunction"
	QualifierKafkaConsumerKeyedHandlerConfiguration = "QualifierKafkaConsumerKeyedHandlerConfiguration"
	QualifierConsumerFunction                       = "QualifierFlowStateless"
	QualifierConsumerTopic                          = "QualifierConsumerTopic"
)

func RegisterConsumer(
	name string,
	broker string,
	configs []runtime.Configuration[*runtime_sarama.Consumer],
) {
	RegisterConsumerConfig(
		runtime_sarama.WithConsumerBroker(broker),
		runtime_sarama.WithConsumerGroupName(name),
	)
	RegisterConsumerConfig(configs...)
	inverse.Register(QualifierKafkaConsumerConfiguration, ResolveConsumerHandlerConfiguration)
	inverse.Register(QualifierKafkaConsumerConfiguration, ResolveConsumerTopics)
	inverse.RegisterWithConfigurationRequired[*runtime_sarama.Consumer](QualifierKafkaConsumer, QualifierKafkaConsumerConfiguration, runtime_sarama.NewConsumer)

	inverse.RegisterConfiguration[*runtime_sarama.KeyedHandler](QualifierKafkaConsumerKeyedHandlerConfiguration, runtime_sarama.WithKeyedHandlerPrometheus())
	inverse.Register(QualifierKafkaConsumerKeyedHandlerConfiguration, ResolveConsumerKeyFunction)
	inverse.Register(QualifierKafkaConsumerKeyedHandlerConfiguration, ResolveConsumerBatchFunction)
	inverse.RegisterWithConfigurationRequired[*runtime_sarama.KeyedHandler](QualifierKafkaConsumerHandler, QualifierKafkaConsumerKeyedHandlerConfiguration, runtime_sarama.NewKeyedHandler)

	RegisterRuntime(QualifierKafkaConsumer)
}

func RegisterConsumerConfig(config ...runtime.Configuration[*runtime_sarama.Consumer]) {
	inverse.RegisterInstances(QualifierKafkaConsumerConfiguration, config)
}

func ResolveConsumerHandlerConfiguration(ctx context.Context) (runtime.Configuration[*runtime_sarama.Consumer], error) {
	consumerHandler, getConsumerHandlerError := inverse.GetLast[runtime_sarama.ConsumerHandler](ctx, QualifierKafkaConsumerHandler)
	if getConsumerHandlerError != nil {
		return nil, getConsumerHandlerError
	}

	return runtime_sarama.WithConsumerLoop(consumerHandler), nil
}

// ===================================

type ConsumerFunction struct {
	Topic string
	Fn    stateless.BatchFunction
	Key   stateful.PersistenceIdFunction[structure.Bytes, structure.Bytes]
}

func RegisterConsumerFunctionInstance(instance ConsumerFunction) {
	inverse.RegisterInstance(QualifierConsumerFunction, instance)
}

func RegisterConsumerFunctionInjector(injector inverse.Injector[ConsumerFunction]) {
	inverse.Register(QualifierConsumerFunction, injector)
}

func ResolveConsumerBatchFunction(ctx context.Context) (runtime.Configuration[*runtime_sarama.KeyedHandler], error) {
	allStatelessFunctions, getAllErr := inverse.GetAll[ConsumerFunction](ctx, QualifierConsumerFunction)
	if getAllErr != nil {
		return nil, getAllErr
	}

	if len(allStatelessFunctions) == 0 {
		return nil, errors.New("missing stateless functions")
	}

	retry, retryErr := GetRetry(ctx)
	if retryErr != nil {
		return nil, retryErr
	}

	producer, producerErr := GetKafkaProducer(ctx)
	if producerErr != nil {
		return nil, producerErr
	}

	var batchFn stateless.BatchFunction

	if len(allStatelessFunctions) > 1 {
		topicSwitchConfigs := []runtime.Configuration[*stateless.TopicSwitch]{}
		for _, instance := range allStatelessFunctions {
			topicSwitchConfigs = append(topicSwitchConfigs, stateless.WithTopicSwitchFunction(instance.Topic, instance.Fn))
		}
		batchFn = stateless.NewTopicSwitch(topicSwitchConfigs...,
		)
	} else {
		batchFn = allStatelessFunctions[0].Fn
	}

	batchFn = stateless.NewProducerBatchFunction(
		stateless.WithBatchProducerNextFunction(batchFn),
		stateless.WithBatchProducerRuntime(producer),
		stateless.WithBatchProducerPrometheus(),
	)

	batchFn = stateless.NewBatchRetry(
		stateless.WithBatchRetryNextFunction(batchFn),
		stateless.WithBatchRetryRuntime(retry),
		stateless.WithBatchRetryPrometheus(),
	)

	return runtime_sarama.WithKeyedHandlerFunction(batchFn), nil
}

func ResolveConsumerKeyFunction(ctx context.Context) (runtime.Configuration[*runtime_sarama.KeyedHandler], error) {
	allStatelessFunctions, getAllErr := inverse.GetAll[ConsumerFunction](ctx, QualifierConsumerFunction)
	if getAllErr != nil {
		return nil, getAllErr
	}

	if len(allStatelessFunctions) == 0 {
		return nil, errors.New("missing stateless functions")
	}

	var keyFn stateful.PersistenceIdFunction[structure.Bytes, structure.Bytes]

	if len(allStatelessFunctions) > 1 {
		keySwitchConfigs := []runtime.Configuration[*stateful.PersistenceIdSwitch]{}
		for _, instance := range allStatelessFunctions {
			keySwitchConfigs = append(keySwitchConfigs, stateful.WithPersistenceIdSwitchPersistenceIdFunction(instance.Topic, instance.Key))
		}
		keyFn = stateful.NewPersistenceIdSwitch(keySwitchConfigs...)
	} else {
		keyFn = allStatelessFunctions[0].Key
	}

	return runtime_sarama.WithKeyedHandlerKeyFunction(keyFn), nil
}

func ResolveConsumerTopics(ctx context.Context) (runtime.Configuration[*runtime_sarama.Consumer], error) {
	allStatelessFunctions, getAllErr := inverse.GetAll[ConsumerFunction](ctx, QualifierConsumerFunction)
	if getAllErr != nil {
		return nil, getAllErr
	}

	if len(allStatelessFunctions) == 0 {
		return nil, errors.New("missing stateless functions")
	}

	topics := []string{}

	for _, instance := range allStatelessFunctions {
		topics = append(topics, instance.Topic)
	}

	return runtime_sarama.WithConsumerTopic(topics...), nil
}
