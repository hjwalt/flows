package flows

import (
	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateless"
)

func KafkaProducer(ctrl runtime.Controller, configs []runtime.Configuration[*runtime_sarama.Producer]) runtime.Producer {
	producerConfig := append(
		configs,
		runtime_sarama.WithProducerRuntimeController(ctrl),
	)
	return runtime_sarama.NewProducer(producerConfig...)
}

func KafkaConsumerSingle(ctrl runtime.Controller, fn stateless.SingleFunction, configs []runtime.Configuration[*runtime_sarama.Consumer]) runtime.Runtime {

	// sarama consumer loop
	consumerLoop := runtime_sarama.NewSingleLoop(
		runtime_sarama.WithLoopSingleFunction(fn),
		runtime_sarama.WithLoopSinglePrometheus(),
	)

	// consumer runtime
	consumerConfig := append(
		configs,
		runtime_sarama.WithConsumerRuntimeController(ctrl),
		runtime_sarama.WithConsumerLoop(consumerLoop),
	)
	return runtime_sarama.NewConsumer(consumerConfig...)
}
