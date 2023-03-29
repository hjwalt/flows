package example

import (
	"context"

	"github.com/hjwalt/flows"
	"github.com/hjwalt/flows/format"
	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/logger"
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

func WordRemapRun() error {
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
			runtime_sarama.WithConsumerGroupName("test"),
		},
	}

	multi := statelessFunctionConfiguration.Runtime()
	return multi.Start()
}
