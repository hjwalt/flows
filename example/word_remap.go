package example

import (
	"context"

	"github.com/hjwalt/flows"
	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/runtime"
	"go.uber.org/zap"
)

func WordRemapStatelessFunction(c context.Context, m flow.Message[string, string]) (*flow.Message[string, string], error) {
	logger.Info("count", zap.String("remap", m.Value+" updated"), zap.String("key", m.Key))
	return &flow.Message[string, string]{
		Topic:   "word-updated",
		Key:     m.Key,
		Value:   m.Value + " updated",
		Headers: m.Headers,
	}, nil
}

func WordRemap() runtime.Runtime {
	statelessFunctionConfiguration := flows.StatelessOneToOneConfiguration[string, string, string, string]{
		Name:         "flows-word-remap",
		InputTopic:   flow.StringTopic("word"),
		OutputTopic:  flow.StringTopic("word-updated"),
		Function:     WordRemapStatelessFunction,
		InputBroker:  "localhost:9092",
		OutputBroker: "localhost:9092",
		HttpPort:     8081,
	}

	return statelessFunctionConfiguration.Runtime()
}
