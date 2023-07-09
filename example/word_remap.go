package example

import (
	"context"

	"github.com/hjwalt/flows"
	"github.com/hjwalt/flows/format"
	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/runtime"
)

func WordRemapStatelessFunction(c context.Context, m message.Message[string, string]) (*message.Message[string, string], error) {
	logger.Info("applying")
	return &message.Message[string, string]{
		Topic:   "word-updated",
		Key:     m.Key,
		Value:   m.Value + " updated",
		Headers: m.Headers,
	}, nil
}

func WordRemap() runtime.Runtime {
	statelessFunctionConfiguration := flows.StatelessSingleFunctionConfiguration{
		StatelessFunction: stateless.ConvertOneToOne(
			WordRemapStatelessFunction,
			format.String(),
			format.String(),
			format.String(),
			format.String(),
		),
		KafkaProducerConfiguration: []runtime.Configuration[*runtime_sarama.Producer]{
			runtime_sarama.WithProducerSaramaConfig(runtime_sarama.DefaultConfiguration()),
			runtime_sarama.WithProducerBroker("localhost:9092"),
		},
		KafkaConsumerConfiguration: []runtime.Configuration[*runtime_sarama.Consumer]{
			runtime_sarama.WithConsumerSaramaConfig(runtime_sarama.DefaultConfiguration()),
			runtime_sarama.WithConsumerBroker("localhost:9092"),
			runtime_sarama.WithConsumerTopic("word-updated"),
			runtime_sarama.WithConsumerGroupName("flows-word-remap"),
		},
	}

	return statelessFunctionConfiguration.Runtime()
}
