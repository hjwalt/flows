package example

import (
	"context"

	"github.com/hjwalt/flows"
	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/reflect"
	"github.com/hjwalt/runway/runtime"
	"go.uber.org/zap"
)

func WordJoinPersistenceId(ctx context.Context, m flow.Message[string, string]) (string, error) {
	return m.Key, nil
}

func WordJoinCountFunction(c context.Context, m flow.Message[string, string], s stateful.State[*WordJoinState]) (*flow.Message[string, string], stateful.State[*WordJoinState], error) {
	logger.Info("applying")

	// setting defaults
	if s.Content == nil {
		s.Content = &WordJoinState{Count: 0, Word: ""}
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

func WordJoinWordFunction(c context.Context, m flow.Message[string, string], s stateful.State[*WordJoinState]) (*flow.Message[string, string], stateful.State[*WordJoinState], error) {
	logger.Info("applying")

	// setting defaults
	if s.Content == nil {
		s.Content = &WordJoinState{Count: 0, Word: ""}
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

func WordJoin() runtime.Runtime {
	joinFunctionConfiguration := flows.JoinPostgresqlFunctionConfiguration{
		StatefulFunctions: map[string]stateful.SingleFunction{
			"word": stateful.ConvertOneToOne(
				WordJoinCountFunction,
				format.Protobuf[*WordJoinState](),
				format.String(),
				format.String(),
				format.String(),
				format.String(),
			),
			"word-type": stateful.ConvertOneToOne(
				WordJoinWordFunction,
				format.Protobuf[*WordJoinState](),
				format.String(),
				format.String(),
				format.String(),
				format.String(),
			),
		},
		PersistenceIdFunctions: map[string]stateful.PersistenceIdFunction[[]byte, []byte]{
			"word": stateful.ConvertPersistenceId(
				WordJoinPersistenceId,
				format.String(),
				format.String(),
			),
			"word-type": stateful.ConvertPersistenceId(
				WordJoinPersistenceId,
				format.String(),
				format.String(),
			),
		},
		Name:                     "flows-word-join",
		InputBroker:              "localhost:9092",
		OutputBroker:             "localhost:9092",
		IntermediateTopicName:    "word-join-intermediate",
		PostgresTable:            "public.flows_state",
		PostgresConnectionString: "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
		HttpPort:                 8081,
	}

	return joinFunctionConfiguration.Runtime()
}
