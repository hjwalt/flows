package example

import (
	"context"
	"net/http"

	"github.com/hjwalt/flows"
	"github.com/hjwalt/flows/format"
	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/protobuf"
	"github.com/hjwalt/flows/router"
	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/stateful_bun"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/reflect"
	"github.com/uptrace/bunrouter"
	"go.uber.org/zap"
)

func WordCountStatefulFunction(c context.Context, m message.Message[message.Bytes, message.Bytes], inState stateful.SingleState[message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], stateful.SingleState[message.Bytes], error) {
	logger.Info("applying")

	// format conversion to something more usable, will be abstracted out also in the future
	protoState, protoStateMapErr := stateful.ConvertSingleState(inState, format.Bytes(), format.Protobuf[*WordCountState]())
	if protoStateMapErr != nil {
		return make([]message.Message[[]byte, []byte], 0), inState, protoStateMapErr
	}

	// setting defaults
	if protoState.Content == nil {
		protoState.Content = &WordCountState{Count: 0}
	}

	// update state
	protoState.Content.Count += 1

	logger.Info("count", zap.Int64("count", protoState.Content.Count))

	// map back to bytes
	nextByteState, nextByteStateMapErr := stateful.ConvertSingleState(protoState, format.Protobuf[*WordCountState](), format.Bytes())
	if nextByteStateMapErr != nil {
		return make([]message.Message[[]byte, []byte], 0), inState, nextByteStateMapErr
	}

	// create output message
	outMessage := message.Message[message.Bytes, message.Bytes]{
		Topic: "word-count",
		Key:   m.Key,
		Value: []byte(reflect.GetString(protoState.Content.Count)),
	}

	return []message.Message[[]byte, []byte]{outMessage}, nextByteState, nil
}

func WordCountPersistenceId(ctx context.Context, m message.Message[message.Bytes, message.Bytes]) (string, error) {
	return string(m.Key), nil
}

func WordCountRun() error {
	statefulFunctionConfiguration := flows.StatefulPostgresqlFunctionConfiguration{
		StatefulFunction:      WordCountStatefulFunction,
		PersistenceIdFunction: WordCountPersistenceId,
		PersistenceTableName:  "public.flows_state",

		PostgresqlConfiguration: []runtime.Configuration[*stateful_bun.PostgresqlConnection]{
			stateful_bun.WithApplicationName("flows"),
			stateful_bun.WithConnectionString("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
		},
		KafkaProducerConfiguration: []runtime.Configuration[*runtime_sarama.Producer]{
			runtime_sarama.WithProducerSaramaConfig(runtime_sarama.DefaultConfiguration()),
			runtime_sarama.WithProducerBroker("localhost:9092"),
		},
		KafkaConsumerConfiguration: []runtime.Configuration[*runtime_sarama.Consumer]{
			runtime_sarama.WithConsumerSaramaConfig(runtime_sarama.DefaultConfiguration()),
			runtime_sarama.WithConsumerBroker("localhost:9092"),
			runtime_sarama.WithConsumerTopic("word"),
			runtime_sarama.WithConsumerGroupName("test"),
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
		},
	}

	multi := statefulFunctionConfiguration.Runtime()
	return multi.Start()
}

type TestResponse struct {
	Message string
}
