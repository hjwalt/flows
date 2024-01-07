package example_word_count

import (
	"context"

	"github.com/avast/retry-go"
	"github.com/hjwalt/flows"
	"github.com/hjwalt/flows/example"
	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/runtime_neo4j"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/reflect"
	"github.com/hjwalt/runway/runtime"
	"go.uber.org/zap"
)

const (
	Instance = "flows-word-count"
)

func key(ctx context.Context, m flow.Message[string, string]) (string, error) {
	return m.Key, nil
}

func fn(c context.Context, m flow.Message[string, string], s stateful.State[*example.WordCountState]) (*flow.Message[string, string], stateful.State[*example.WordCountState], error) {
	// setting defaults
	if s.Content == nil {
		s.Content = &example.WordCountState{Count: 0}
	}

	// update state
	s.Content.Count += 1

	logger.Info("count", zap.Int64("count", s.Content.Count), zap.String("key", m.Key))

	// create output message
	outMessage := flow.Message[string, string]{
		Topic: "word-count",
		Key:   m.Key,
		Value: reflect.GetString(s.Content.Count),
	}

	return &outMessage, s, nil
}

func Registrar(ci inverse.Container) flows.Prebuilt {
	return flows.StatefulPostgresqlOneToOneFunctionConfiguration[*example.WordCountState, string, string, string, string]{
		Name:                     Instance,
		InputTopic:               flow.StringTopic("word"),
		OutputTopic:              flow.StringTopic("word-count"),
		Function:                 fn,
		InputBroker:              "localhost:9092",
		OutputBroker:             "localhost:9092",
		HttpPort:                 8081,
		StateFormat:              format.Protobuf[*example.WordCountState](),
		StateKeyFunction:         key,
		PostgresTable:            "public.flows_count_state",
		PostgresConnectionString: "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",

		Neo4jConfiguration: []runtime.Configuration[*runtime_neo4j.Neo4JConnectionBasicAuth]{
			runtime_neo4j.WithPass("localhost"),
			runtime_neo4j.WithUser("neo4j"),
			runtime_neo4j.WithUrl("neo4j://localhost:7687"),
		},
		// Optional configurations
		RetryConfiguration: []runtime.Configuration[*runtime_retry.Retry]{
			runtime_retry.WithRetryOption(
				retry.Attempts(3),
			),
			runtime_retry.WithAbsorbError(true),
		},
	}
}

type TestResponse struct {
	Message string
}

func Register(m flows.Main) {
	m.Prebuilt(Instance, Registrar)
}
