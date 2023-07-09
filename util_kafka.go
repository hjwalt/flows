package flows

import (
	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/runtime"
)

func KafkaProducer(configs []runtime.Configuration[*runtime_sarama.Producer]) message.Producer {
	return runtime_sarama.NewProducer(configs...)
}

func KafkaConsumerSingle(fn stateless.SingleFunction, configs []runtime.Configuration[*runtime_sarama.Consumer]) runtime.Runtime {

	// sarama consumer loop
	consumerLoop := runtime_sarama.NewSingleLoop(
		runtime_sarama.WithLoopSingleFunction(fn),
		runtime_sarama.WithLoopSinglePrometheus(),
	)

	// consumer runtime
	consumerConfig := append(
		configs,
		runtime_sarama.WithConsumerLoop(consumerLoop),
	)
	return runtime_sarama.NewConsumer(consumerConfig...)
}
