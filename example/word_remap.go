package example

import (
	"context"

	"github.com/hjwalt/flows"
	"github.com/hjwalt/flows/format"
	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/runway/logger"
)

func WordRemapStatelessFunction(c context.Context, m message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {
	logger.Info("applying")

	// format map to proto
	stringMessage, mapErr := message.Convert(m, format.Bytes(), format.Bytes(), format.String(), format.String())
	if mapErr != nil {
		return make([]message.Message[[]byte, []byte], 0), mapErr
	}

	stringMessage.Value = stringMessage.Value + " updated"
	stringMessage.Topic = "word-updated"

	// format map from proto
	outMessage, outMessageMapErr := message.Convert(stringMessage, format.String(), format.String(), format.Bytes(), format.Bytes())
	if outMessageMapErr != nil {
		return make([]message.Message[[]byte, []byte], 0), outMessageMapErr
	}

	return []message.Message[[]byte, []byte]{outMessage}, nil
}

func WordRemapRun() error {
	statelessFunctionConfiguration := flows.StatelessSingleFunctionConfiguration{
		StatelessFunction: WordRemapStatelessFunction,
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
