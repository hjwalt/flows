package flows

import (
	"context"
	"time"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
)

const (
	QualifierKafkaProducerConfiguration  = "QualifierKafkaProducerConfiguration"
	QualifierKafkaProducer               = "QualifierKafkaProducer"
	QualifierKafkaConsumerConfiguration  = "QualifierKafkaConsumerConfiguration"
	QualifierKafkaConsumer               = "QualifierKafkaConsumer"
	QualifierKafkaConsumerHandler        = "QualifierKafkaConsumerHandler"
	QualifierKafkaConsumerSingleFunction = "QualifierKafkaConsumerSingleFunction"
	QualifierKafkaConsumerBatchFunction  = "QualifierKafkaConsumerBatchFunction"
	QualifierKafkaConsumerKeyFunction    = "QualifierKafkaConsumerKeyFunction"
)

// Producer
func RegisterProducerConfig(config []runtime.Configuration[*runtime_sarama.Producer]) {
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

func GetKafkaProducer(ctx context.Context) (message.Producer, error) {
	return inverse.GetLast[message.Producer](ctx, QualifierKafkaProducer)
}

// Consumer
func RegisterConsumer() {
	inverse.RegisterWithConfigurationRequired[*runtime_sarama.Consumer](
		QualifierKafkaConsumer,
		QualifierKafkaConsumerConfiguration,
		runtime_sarama.NewConsumer,
	)
	inverse.Register(QualifierRuntime, InjectorRuntime(QualifierKafkaConsumer))
	inverse.Register[runtime.Configuration[*runtime_sarama.Consumer]](QualifierKafkaConsumerConfiguration, InjectorKafkaConsumerHandlerConfiguration)
}

func InjectorKafkaConsumerHandlerConfiguration(ctx context.Context) (runtime.Configuration[*runtime_sarama.Consumer], error) {
	consumerHandler, getConsumerHandlerError := inverse.GetLast[runtime_sarama.ConsumerLoop](ctx, QualifierKafkaConsumerHandler)
	if getConsumerHandlerError != nil {
		return nil, getConsumerHandlerError
	}

	return runtime_sarama.WithConsumerLoop(consumerHandler), nil
}

func RegisterConsumerKeyedConfig(config []runtime.Configuration[*runtime_sarama.Consumer]) {
	inverse.RegisterInstances(QualifierKafkaConsumerConfiguration, config)
	inverse.Register[runtime_sarama.ConsumerLoop](QualifierKafkaConsumerHandler, InjectorKafkaConsumerKeyedHandler)
}

func InjectorKafkaConsumerKeyedHandler(ctx context.Context) (runtime_sarama.ConsumerLoop, error) {
	batchFunction, getBatchFunctionError := inverse.GetLast[stateless.BatchFunction](ctx, QualifierKafkaConsumerBatchFunction)
	if getBatchFunctionError != nil {
		return nil, getBatchFunctionError
	}

	keyFunction, getKeyFunctionError := inverse.GetLast[stateful.PersistenceIdFunction[message.Bytes, message.Bytes]](ctx, QualifierKafkaConsumerKeyFunction)
	if getKeyFunctionError != nil {
		return nil, getKeyFunctionError
	}

	// sarama consumer loop
	consumerLoop := runtime_sarama.NewKeyedHandler(
		runtime_sarama.WithKeyedHandlerMaxBufferred(1000),
		runtime_sarama.WithKeyedHandlerMaxDelay(100*time.Millisecond),
		runtime_sarama.WithKeyedHandlerFunction(batchFunction),
		runtime_sarama.WithKeyedHandlerKeyFunction(keyFunction),
		runtime_sarama.WithKeyedHandlerPrometheus(),
	)

	return consumerLoop, nil
}
