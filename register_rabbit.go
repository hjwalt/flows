package flows

import (
	"context"

	"github.com/hjwalt/flows/runtime_rabbit"
	"github.com/hjwalt/flows/task"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
	"github.com/hjwalt/runway/structure"
)

const (
	QualifierRabbitProducer         = "QualifierRabbitProducer"
	QualifierRabbitConsumerExecutor = "QualifierRabbitConsumerExecutor"
	QualifierRabbitConsumer         = "QualifierRabbitConsumer"
)

// Producer

func RegisterRabbitProducer(
	container inverse.Container,
	name string,
	rabbitConnectionString string,
	configs []runtime.Configuration[*runtime_rabbit.Producer],
) {

	resolver := runtime.NewResolver[*runtime_rabbit.Producer, task.Producer](
		QualifierRabbitProducer,
		container,
		true,
		runtime_rabbit.NewProducer,
	)

	resolver.AddConfigVal(runtime_rabbit.WithProducerName(name))
	resolver.AddConfigVal(runtime_rabbit.WithProducerConnectionString(rabbitConnectionString))

	for _, config := range configs {
		resolver.AddConfigVal(config)
	}

	resolver.Register()

	RegisterRuntime(QualifierRabbitProducer, container)
}

func GetRabbitProducer(ctx context.Context, ci inverse.Container) (task.Producer, error) {
	return inverse.GenericGetLast[task.Producer](ci, ctx, QualifierRabbitProducer)
}

// Consumer

func RegisterRabbitConsumer(
	container inverse.Container,
	name string,
	rabbitConnectionString string,
	channel string,
	configs []runtime.Configuration[*runtime_rabbit.Consumer],
) {

	resolver := runtime.NewResolver[*runtime_rabbit.Consumer, runtime.Runtime](
		QualifierRabbitConsumer,
		container,
		true,
		runtime_rabbit.NewConsumer,
	)

	resolver.AddConfigVal(runtime_rabbit.WithConsumerName(name))
	resolver.AddConfigVal(runtime_rabbit.WithConsumerQueueName(channel))
	resolver.AddConfigVal(runtime_rabbit.WithConsumerConnectionString(rabbitConnectionString))
	resolver.AddConfig(InjectorRabbitConsumerExecutorConfiguration)

	for _, config := range configs {
		resolver.AddConfigVal(config)
	}

	resolver.Register()

	RegisterRuntime(QualifierRabbitConsumer, container)
}

func RegisterRabbitConsumerExecutor(ci inverse.Container, injector func(ctx context.Context, ci inverse.Container) (task.Executor[structure.Bytes], error)) {
	inverse.GenericAdd(ci, QualifierRabbitConsumerExecutor, injector)
}

func InjectorRabbitConsumerExecutorConfiguration(ctx context.Context, ci inverse.Container) (runtime.Configuration[*runtime_rabbit.Consumer], error) {
	handler, getHandlerError := inverse.GenericGetLast[task.Executor[structure.Bytes]](ci, ctx, QualifierRabbitConsumerExecutor)
	if getHandlerError != nil {
		return nil, getHandlerError
	}
	return runtime_rabbit.WithConsumerHandler(handler), nil
}

// ===================================

func RegisterRabbitProducerConfig(ci inverse.Container, configs ...runtime.Configuration[*runtime_rabbit.Producer]) {
	for _, config := range configs {
		ci.AddVal(runtime.QualifierConfig(QualifierRabbitProducer), config)
	}
}

func RegisterRabbitConsumerConfig(ci inverse.Container, configs ...runtime.Configuration[*runtime_rabbit.Consumer]) {
	for _, config := range configs {
		ci.AddVal(runtime.QualifierConfig(QualifierRabbitConsumer), config)
	}
}
