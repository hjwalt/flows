package flows

import (
	"context"
	"errors"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
)

const (
	QualifierKafkaProducerConfiguration       = "QualifierKafkaProducerConfiguration"
	QualifierKafkaProducer                    = "QualifierKafkaProducer"
	QualifierKafkaConsumerSingleConfiguration = "QualifierKafkaConsumerSingleConfiguration"
	QualifierKafkaConsumerSingle              = "QualifierKafkaConsumerSingle"
	QualifierKafkaConsumerSingleFunction      = "QualifierKafkaConsumerSingleFunction"
)

// Producer
func RegisterProducerConfig(config []runtime.Configuration[*runtime_sarama.Producer]) {
	inverse.RegisterInstances(QualifierKafkaProducerConfiguration, config)
}

func RegisterProducer() {
	inverse.Register(QualifierKafkaProducer, InjectorKafkaProducer)
	inverse.Register(QualifierRuntime, InjectorRuntime(QualifierKafkaProducer))
}

func InjectorKafkaProducer(ctx context.Context) (message.Producer, error) {
	configurations, getConfigurationError := inverse.GetAll[runtime.Configuration[*runtime_sarama.Producer]](ctx, QualifierKafkaProducerConfiguration)
	if getConfigurationError != nil {
		return nil, getConfigurationError
	}
	return runtime_sarama.NewProducer(configurations...), nil
}

func GetKafkaProducer(ctx context.Context) (message.Producer, error) {
	return inverse.GetLast[message.Producer](ctx, QualifierKafkaProducer)
}

// Consumer

func RegisterConsumerSingleConfig(config []runtime.Configuration[*runtime_sarama.Consumer]) {
	inverse.RegisterInstances(QualifierKafkaConsumerSingleConfiguration, config)
}

func RegisterConsumerSingle() {
	inverse.Register(QualifierKafkaConsumerSingle, InjectorKafkaConsumerSingle)
	inverse.Register(QualifierRuntime, InjectorRuntime(QualifierKafkaConsumerSingle))
}

func InjectorKafkaConsumerSingle(ctx context.Context) (runtime.Runtime, error) {
	configurations, getConfigurationError := inverse.GetAll[runtime.Configuration[*runtime_sarama.Consumer]](ctx, QualifierKafkaConsumerSingleConfiguration)
	if getConfigurationError != nil && !errors.Is(getConfigurationError, inverse.ErrNotInjected) {
		return nil, getConfigurationError
	}

	singleFunction, getSingleFunctionError := inverse.GetLast[stateless.SingleFunction](ctx, QualifierKafkaConsumerSingleFunction)
	if getSingleFunctionError != nil {
		return nil, getSingleFunctionError
	}

	// sarama consumer loop
	consumerLoop := runtime_sarama.NewSingleLoop(
		runtime_sarama.WithLoopSingleFunction(singleFunction),
		runtime_sarama.WithLoopSinglePrometheus(),
	)

	// consumer runtime
	consumerConfig := append(
		configurations,
		runtime_sarama.WithConsumerLoop(consumerLoop),
	)

	return runtime_sarama.NewConsumer(consumerConfig...), nil
}
