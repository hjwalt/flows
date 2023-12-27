package example_word_remap

import (
	"context"

	"github.com/hjwalt/flows"
	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/runtime_neo4j"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/runtime"
	"go.uber.org/zap"
)

const (
	Instance = "flows-word-remap"
)

func WordRemapStatelessFunction(c context.Context, m flow.Message[string, string]) (*flow.Message[string, string], error) {
	logger.Info("remap", zap.String("remap", m.Value+" updated"), zap.String("key", m.Key))
	return &flow.Message[string, string]{
		Topic:   "word-updated",
		Key:     m.Key,
		Value:   m.Value + " updated",
		Headers: m.Headers,
	}, nil
}

func Registrar(ci inverse.Container) flows.Prebuilt {
	return flows.StatelessOneToOneConfiguration[string, string, string, string]{
		Name:         Instance,
		InputTopic:   flow.StringTopic("word"),
		OutputTopic:  flow.StringTopic("word-updated"),
		Function:     WordRemapStatelessFunction,
		InputBroker:  "localhost:9092",
		OutputBroker: "localhost:9092",
		HttpPort:     8081,
		Neo4jConfiguration: []runtime.Configuration[*runtime_neo4j.Neo4JConnectionBasicAuth]{
			runtime_neo4j.WithPass("localhost"),
			runtime_neo4j.WithUser("neo4j"),
			runtime_neo4j.WithUrl("neo4j://localhost:7687"),
		},
	}
}

func Register(m flows.Main) {
	m.Prebuilt(Instance, Registrar)
}
