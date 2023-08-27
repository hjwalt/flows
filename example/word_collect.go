package example

import (
	"context"
	"time"

	"github.com/hjwalt/flows"
	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/reflect"
	"github.com/hjwalt/runway/runtime"
	"go.uber.org/zap"
)

func WordCollectPersistenceId(ctx context.Context, m flow.Message[string, string]) (string, error) {
	return m.Key, nil
}

func WordCollectAggregator(c context.Context, m flow.Message[string, string], s stateful.State[*WordCollectState]) (stateful.State[*WordCollectState], error) {
	logger.Info("applying")

	// setting defaults
	if s.Content == nil {
		s.Content = &WordCollectState{Count: 0}
	}

	// update state
	s.Content.Word = m.Key
	s.Content.Count += 1

	logger.Info("count", zap.Int64("count", s.Content.Count), zap.String("key", m.Key))

	return s, nil
}

func WordCollectCollector(c context.Context, persistenceId string, s stateful.State[*WordCollectState]) (*flow.Message[string, string], error) {

	// create output message
	outMessage := flow.Message[string, string]{
		Topic: "word-count",
		Key:   s.Content.Word,
		Value: reflect.GetString(s.Content.Count),
	}

	return &outMessage, nil
}

func WordCollect() runtime.Runtime {
	inverse.RegisterConfiguration[*runtime_sarama.KeyedHandler](flows.QualifierKafkaKeyedHandlerConfiguration, runtime_sarama.WithKeyedHandlerMaxPerKey(1000))
	inverse.RegisterConfiguration[*runtime_sarama.KeyedHandler](flows.QualifierKafkaKeyedHandlerConfiguration, runtime_sarama.WithKeyedHandlerMaxBufferred(10000))
	inverse.RegisterConfiguration[*runtime_sarama.KeyedHandler](flows.QualifierKafkaKeyedHandlerConfiguration, runtime_sarama.WithKeyedHandlerMaxDelay(1*time.Second))

	runtimeConfig := flows.CollectorOneToOneConfiguration[*WordCollectState, string, string, string, string]{
		Name:             "flows-word-collect",
		InputTopic:       flow.StringTopic("word"),
		OutputTopic:      flow.StringTopic("word-collect"),
		Aggregator:       WordCollectAggregator,
		Collector:        WordCollectCollector,
		InputBroker:      "localhost:9092",
		OutputBroker:     "localhost:9092",
		HttpPort:         8081,
		StateFormat:      format.Protobuf[*WordCollectState](),
		StateKeyFunction: WordCollectPersistenceId,
	}

	return runtimeConfig.Runtime()
}
