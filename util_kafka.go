package flows

import (
	"context"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
	"github.com/hjwalt/runway/structure"
)

const (
	QualifierKafkaProducerConfiguration     = "QualifierKafkaProducerConfiguration"
	QualifierKafkaProducer                  = "QualifierKafkaProducer"
	QualifierKafkaConsumerConfiguration     = "QualifierKafkaConsumerConfiguration"
	QualifierKafkaConsumer                  = "QualifierKafkaConsumer"
	QualifierKafkaConsumerHandler           = "QualifierKafkaConsumerHandler"
	QualifierKafkaConsumerSingleFunction    = "QualifierKafkaConsumerSingleFunction"
	QualifierKafkaConsumerBatchFunction     = "QualifierKafkaConsumerBatchFunction"
	QualifierKafkaConsumerKeyFunction       = "QualifierKafkaConsumerKeyFunction"
	QualifierKafkaKeyedHandlerConfiguration = "QualifierKafkaKeyedHandlerConfiguration"
)

// Producer
func RegisterProducerConfig(config ...runtime.Configuration[*runtime_sarama.Producer]) {
	inverse.RegisterInstances(QualifierKafkaProducerConfiguration, config)
}

func RegisterProducer() {
	inverse.RegisterWithConfigurationRequired[*runtime_sarama.Producer](
		QualifierKafkaProducer,
		QualifierKafkaProducerConfiguration,
		runtime_sarama.NewProducer,
	)
	inverse.Register(QualifierRuntime, InjectorRuntime(QualifierKafkaProducer))
}

func GetKafkaProducer(ctx context.Context) (flow.Producer, error) {
	return inverse.GetLast[flow.Producer](ctx, QualifierKafkaProducer)
}

func RegisterConsumerConfig(config ...runtime.Configuration[*runtime_sarama.Consumer]) {
	inverse.RegisterInstances(QualifierKafkaConsumerConfiguration, config)
}

// Consumer
func RegisterConsumer() {
	inverse.RegisterWithConfigurationRequired[*runtime_sarama.Consumer](
		QualifierKafkaConsumer,
		QualifierKafkaConsumerConfiguration,
		runtime_sarama.NewConsumer,
	)
	inverse.Register(QualifierRuntime, InjectorRuntime(QualifierKafkaConsumer))
	inverse.Register[runtime.Configuration[*runtime_sarama.Consumer]](QualifierKafkaConsumerConfiguration, InjectorConsumerHandlerConfiguration)
}

func InjectorConsumerHandlerConfiguration(ctx context.Context) (runtime.Configuration[*runtime_sarama.Consumer], error) {
	consumerHandler, getConsumerHandlerError := inverse.GetLast[runtime_sarama.ConsumerHandler](ctx, QualifierKafkaConsumerHandler)
	if getConsumerHandlerError != nil {
		return nil, getConsumerHandlerError
	}

	return runtime_sarama.WithConsumerLoop(consumerHandler), nil
}

func RegisterConsumerKeyedConfig() {
	inverse.RegisterWithConfigurationRequired[*runtime_sarama.KeyedHandler](
		QualifierKafkaConsumerHandler,
		QualifierKafkaKeyedHandlerConfiguration,
		runtime_sarama.NewKeyedHandler,
	)
	inverse.Register[runtime.Configuration[*runtime_sarama.KeyedHandler]](QualifierKafkaKeyedHandlerConfiguration, InjectorConsumerKeyedHandlerBatchFunctionConfiguration)
	inverse.Register[runtime.Configuration[*runtime_sarama.KeyedHandler]](QualifierKafkaKeyedHandlerConfiguration, InjectorConsumerKeyedHandlerKeyFunctionConfiguration)
	inverse.RegisterConfiguration[*runtime_sarama.KeyedHandler](QualifierKafkaKeyedHandlerConfiguration, runtime_sarama.WithKeyedHandlerPrometheus())
}

func InjectorConsumerKeyedHandlerBatchFunctionConfiguration(ctx context.Context) (runtime.Configuration[*runtime_sarama.KeyedHandler], error) {
	batchFunction, getBatchFunctionError := inverse.GetLast[stateless.BatchFunction](ctx, QualifierKafkaConsumerBatchFunction)
	if getBatchFunctionError != nil {
		return nil, getBatchFunctionError
	}
	return runtime_sarama.WithKeyedHandlerFunction(batchFunction), nil
}

func InjectorConsumerKeyedHandlerKeyFunctionConfiguration(ctx context.Context) (runtime.Configuration[*runtime_sarama.KeyedHandler], error) {
	keyFunction, getKeyFunctionError := inverse.GetLast[stateful.PersistenceIdFunction[structure.Bytes, structure.Bytes]](ctx, QualifierKafkaConsumerKeyFunction)
	if getKeyFunctionError != nil {
		return nil, getKeyFunctionError
	}
	return runtime_sarama.WithKeyedHandlerKeyFunction(keyFunction), nil
}

func RegisterConsumerKeyedKeyFunction(persistenceIdFunction stateful.PersistenceIdFunction[structure.Bytes, structure.Bytes]) {
	inverse.RegisterInstance[stateful.PersistenceIdFunction[structure.Bytes, structure.Bytes]](QualifierKafkaConsumerKeyFunction, persistenceIdFunction)
}

func RegisterConsumerKeyedFunction(batchFunctionInjector func(ctx context.Context) (stateless.BatchFunction, error)) {
	inverse.Register[stateless.BatchFunction](QualifierKafkaConsumerBatchFunction, batchFunctionInjector)
}
