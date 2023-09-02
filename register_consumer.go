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
	QualifierKafkaConsumer              = "QualifierKafkaConsumer"
	QualifierKafkaConsumerHandler       = "QualifierKafkaConsumerHandler"
	QualifierKafkaConsumerBatchFunction = "QualifierKafkaConsumerBatchFunction"
	QualifierKafkaConsumerKeyFunction   = "QualifierKafkaConsumerKeyFunction"
	QualifierConsumerFunction           = "QualifierFlowStateless"
)

func RegisterConsumer(
	container inverse.Container,
	name string,
	broker string,
	configs []runtime.Configuration[*runtime_sarama.Consumer],
) {

	resolver := runtime.NewResolver[*runtime_sarama.Consumer, runtime.Runtime](
		QualifierKafkaConsumer,
		container,
		true,
		runtime_sarama.NewConsumer,
	)

	resolver.AddConfigVal(runtime_sarama.WithConsumerBroker(broker))
	resolver.AddConfigVal(runtime_sarama.WithConsumerGroupName(name))
	resolver.AddConfig(ResolveConsumerHandlerConfiguration)
	resolver.AddConfig(ResolveConsumerTopics)

	for _, config := range configs {
		resolver.AddConfigVal(config)
	}

	resolver.Register()

	RegisterRuntime(QualifierKafkaConsumer, container)

	consumerKeyedHandlerResolver := runtime.NewResolver[*runtime_sarama.KeyedHandler, runtime_sarama.ConsumerHandler](
		QualifierKafkaConsumerHandler,
		container,
		true,
		runtime_sarama.NewKeyedHandler,
	)

	consumerKeyedHandlerResolver.AddConfigVal(runtime_sarama.WithKeyedHandlerPrometheus())
	consumerKeyedHandlerResolver.AddConfig(ResolveConsumerKeyFunction)
	consumerKeyedHandlerResolver.AddConfig(ResolveConsumerBatchFunction)

	consumerKeyedHandlerResolver.Register()
}

func ResolveConsumerHandlerConfiguration(ctx context.Context, ci inverse.Container) (runtime.Configuration[*runtime_sarama.Consumer], error) {
	consumerHandler, getConsumerHandlerError := inverse.GenericGetLast[runtime_sarama.ConsumerHandler](ci, ctx, QualifierKafkaConsumerHandler)
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

func RegisterConsumerFunctionInstance(instance ConsumerFunction, ci inverse.Container) {
	ci.AddVal(QualifierConsumerFunction, instance)
}

func RegisterConsumerFunctionInjector(injector inverse.Injector[ConsumerFunction], ci inverse.Container) {
	inverse.GenericAdd(ci, QualifierConsumerFunction, injector)
}

func ResolveConsumerBatchFunction(ctx context.Context, ci inverse.Container) (runtime.Configuration[*runtime_sarama.KeyedHandler], error) {
	allStatelessFunctions, getAllErr := inverse.GenericGetAll[ConsumerFunction](ci, ctx, QualifierConsumerFunction)
	if getAllErr != nil {
		return nil, getAllErr
	}

	if len(allStatelessFunctions) == 0 {
		return nil, errors.New("missing stateless functions")
	}

	retry, retryErr := GetRetry(ctx, ci)
	if retryErr != nil {
		return nil, retryErr
	}

	producer, producerErr := GetKafkaProducer(ctx, ci)
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

func ResolveConsumerKeyFunction(ctx context.Context, ci inverse.Container) (runtime.Configuration[*runtime_sarama.KeyedHandler], error) {
	allStatelessFunctions, getAllErr := inverse.GenericGetAll[ConsumerFunction](ci, ctx, QualifierConsumerFunction)
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

func ResolveConsumerTopics(ctx context.Context, ci inverse.Container) (runtime.Configuration[*runtime_sarama.Consumer], error) {
	allStatelessFunctions, getAllErr := inverse.GenericGetAll[ConsumerFunction](ci, ctx, QualifierConsumerFunction)
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

// ===================================

func RegisterConsumerConfig(ci inverse.Container, configs ...runtime.Configuration[*runtime_sarama.Consumer]) {
	for _, config := range configs {
		ci.AddVal(runtime.QualifierConfig(QualifierKafkaConsumer), config)
	}
}

func RegisterConsumerKeyedHandlerConfig(ci inverse.Container, configs ...runtime.Configuration[*runtime_sarama.KeyedHandler]) {
	for _, config := range configs {
		ci.AddVal(runtime.QualifierConfig(QualifierKafkaConsumerHandler), config)
	}
}
