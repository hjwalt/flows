package example_word_join

import (
	"context"

	"github.com/hjwalt/flows"
	"github.com/hjwalt/flows/example"
	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/reflect"
	"go.uber.org/zap"
)

const (
	Instance = "flows-word-join"
)

func key(ctx context.Context, m flow.Message[string, string]) (string, error) {
	return m.Key, nil
}

func joinCount(c context.Context, m flow.Message[string, string], s stateful.State[*example.WordJoinState]) (*flow.Message[string, string], stateful.State[*example.WordJoinState], error) {
	logger.Info("join count")

	// setting defaults
	if s.Content == nil {
		s.Content = &example.WordJoinState{Count: 0, Word: ""}
	}

	// update state
	s.Content.Count += 1

	logger.Info("info", zap.Int64("count", s.Content.Count), zap.String("word", s.Content.Word))

	// create output message
	outMessage := flow.Message[string, string]{
		Topic: "word-join",
		Key:   m.Key,
		Value: reflect.GetString(s.Content.Count) + " " + s.Content.Word,
	}

	return &outMessage, s, nil
}

func joinWord(c context.Context, m flow.Message[string, string], s stateful.State[*example.WordJoinState]) (*flow.Message[string, string], stateful.State[*example.WordJoinState], error) {
	logger.Info("join word")

	// setting defaults
	if s.Content == nil {
		s.Content = &example.WordJoinState{Count: 0, Word: ""}
	}

	// update state
	s.Content.Word = string(m.Value)

	logger.Info("info", zap.Int64("count", s.Content.Count), zap.String("word", s.Content.Word))

	// create output message
	outMessage := flow.Message[string, string]{
		Topic: "word-join",
		Key:   m.Key,
		Value: reflect.GetString(s.Content.Count) + " " + s.Content.Word,
	}

	return &outMessage, s, nil
}

func Registrar(ci inverse.Container) flows.Prebuilt {
	return flows.JoinPostgresqlFunctionConfiguration{
		StatefulFunctions: map[string]stateful.SingleFunction{
			"word": stateful.ConvertOneToOne(
				joinCount,
				format.Protobuf[*example.WordJoinState](),
				format.String(),
				format.String(),
				format.String(),
				format.String(),
			),
			"word-type": stateful.ConvertOneToOne(
				joinWord,
				format.Protobuf[*example.WordJoinState](),
				format.String(),
				format.String(),
				format.String(),
				format.String(),
			),
		},
		PersistenceIdFunctions: map[string]stateful.PersistenceIdFunction[[]byte, []byte]{
			"word": stateful.ConvertPersistenceId(
				key,
				format.String(),
				format.String(),
			),
			"word-type": stateful.ConvertPersistenceId(
				key,
				format.String(),
				format.String(),
			),
		},
		Name:                     Instance,
		InputBroker:              "localhost:9092",
		OutputBroker:             "localhost:9092",
		IntermediateTopicName:    "word-join-intermediate",
		PostgresTable:            "public.flows_join_state",
		PostgresConnectionString: "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
		HttpPort:                 8081,
	}
}

func Register(m flows.Main) {
	m.Prebuilt(Instance, Registrar)
}
