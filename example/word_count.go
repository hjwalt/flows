package example

import (
	"context"
	"net/http"

	"github.com/avast/retry-go"
	"github.com/hjwalt/flows"
	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/protobuf"
	"github.com/hjwalt/flows/router"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/topic"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/reflect"
	"github.com/hjwalt/runway/runtime"
	"github.com/uptrace/bunrouter"
	"go.uber.org/zap"
)

func WordCountPersistenceId(ctx context.Context, m message.Message[string, string]) (string, error) {
	return m.Key, nil
}

func WordCountStatefulFunction(c context.Context, m message.Message[string, string], s stateful.State[*WordCountState]) (*message.Message[string, string], stateful.State[*WordCountState], error) {
	// setting defaults
	if s.Content == nil {
		s.Content = &WordCountState{Count: 0}
	}

	// update state
	s.Content.Count += 1

	logger.Info("count", zap.Int64("count", s.Content.Count), zap.String("key", m.Key))

	// create output message
	outMessage := message.Message[string, string]{
		Topic: "word-count",
		Key:   m.Key,
		Value: reflect.GetString(s.Content.Count),
	}

	return &outMessage, s, nil
}

func WordCount() runtime.Runtime {
	statefulFunctionConfiguration := flows.StatefulPostgresqlOneToOneFunctionConfiguration[*WordCountState, string, string, string, string]{
		Name:                     "flows-word-count",
		InputTopic:               topic.String("word"),
		OutputTopic:              topic.String("word-count"),
		Function:                 WordCountStatefulFunction,
		InputBroker:              "localhost:9092",
		OutputBroker:             "localhost:9092",
		HttpPort:                 8081,
		StateFormat:              format.Protobuf[*WordCountState](),
		StateKeyFunction:         WordCountPersistenceId,
		PostgresTable:            "public.flows_state",
		PostgresConnectionString: "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",

		// Optional configurations
		RetryConfiguration: []runtime.Configuration[*runtime_retry.Retry]{
			runtime_retry.WithRetryOption(
				retry.Attempts(3),
			),
			runtime_retry.WithAbsorbError(true),
		},
		RouteConfiguration: []runtime.Configuration[*runtime_bunrouter.Router]{
			runtime_bunrouter.WithRouterGroup("/api"),
			runtime_bunrouter.WithRouterBunHandler(runtime_bunrouter.GET, "/dummy", func(w http.ResponseWriter, req bunrouter.Request) error {
				state := &protobuf.State{
					State: &protobuf.State_V1{
						V1: &protobuf.StateV1{
							OffsetProgress: map[int32]int64{
								1: 1,
							},
						},
					},
				}

				return router.WriteJson(w, 200, state, format.Protobuf[*protobuf.State]())
			}),
			runtime_bunrouter.WithRouterBunHandler(runtime_bunrouter.GET, "/test", func(w http.ResponseWriter, req bunrouter.Request) error {
				state := TestResponse{
					Message: "test",
				}
				return router.WriteJson(w, 200, state, format.Json[TestResponse]())
			}),
			runtime_bunrouter.WithRouterProducerHandler(runtime_bunrouter.POST, "/produce", func(ctx context.Context, req message.Message[message.Bytes, message.Bytes]) (*message.Message[message.Bytes, message.Bytes], error) {
				return &message.Message[[]byte, []byte]{
						Topic:   "produce-topic",
						Key:     req.Key,
						Value:   req.Value,
						Headers: req.Headers,
					},
					nil
			}),
		},
	}

	return statefulFunctionConfiguration.Runtime()
}

type TestResponse struct {
	Message string
}
