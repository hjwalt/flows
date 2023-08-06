package flows

import (
	"context"

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
func RegisterConsumerSingleConfig(config []runtime.Configuration[*runtime_sarama.Consumer]) {
	inverse.RegisterInstances(QualifierKafkaConsumerSingleConfiguration, config)
	inverse.Register[runtime.Configuration[*runtime_sarama.Consumer]](QualifierKafkaConsumerSingleConfiguration, InjectorKafkaConsumerSingleLoop)
}

func RegisterConsumerSingle() {
	inverse.RegisterWithConfigurationRequired[*runtime_sarama.Consumer](
		QualifierKafkaConsumerSingle,
		QualifierKafkaConsumerSingleConfiguration,
		runtime_sarama.NewConsumer,
	)
	inverse.Register(QualifierRuntime, InjectorRuntime(QualifierKafkaConsumerSingle))
}

func InjectorKafkaConsumerSingleLoop(ctx context.Context) (runtime.Configuration[*runtime_sarama.Consumer], error) {
	singleFunction, getSingleFunctionError := inverse.GetLast[stateless.SingleFunction](ctx, QualifierKafkaConsumerSingleFunction)
	if getSingleFunctionError != nil {
		return nil, getSingleFunctionError
	}

	// sarama consumer loop
	consumerLoop := runtime_sarama.NewSingleLoop(
		runtime_sarama.WithLoopSingleFunction(singleFunction),
		runtime_sarama.WithLoopSinglePrometheus(),
	)

	return runtime_sarama.WithConsumerLoop(consumerLoop), nil
}
