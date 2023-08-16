package example

import (
	"context"

	"github.com/hjwalt/flows"
	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/topic"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/runtime"
	"go.uber.org/zap"
)

func WordRemapStatelessFunction(c context.Context, m message.Message[string, string]) (*message.Message[string, string], error) {
	logger.Info("count", zap.String("remap", m.Value+" updated"), zap.String("key", m.Key))
	return &message.Message[string, string]{
		Topic:   "word-updated",
		Key:     m.Key,
		Value:   m.Value + " updated",
		Headers: m.Headers,
	}, nil
}

func WordRemap() runtime.Runtime {
	statelessFunctionConfiguration := flows.StatelessOneToOneConfiguration[string, string, string, string]{
		Name:         "flows-word-remap",
		InputTopic:   topic.String("word"),
		OutputTopic:  topic.String("word-updated"),
		Function:     WordRemapStatelessFunction,
		InputBroker:  "localhost:9092",
		OutputBroker: "localhost:9092",
		HttpPort:     8081,
	}

	return statelessFunctionConfiguration.Runtime()
}
