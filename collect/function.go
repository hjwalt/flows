package collect

import (
	"context"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/stateful"
)

type Aggregator[S any, K any, V any] func(context.Context, message.Message[K, V], stateful.State[S]) (stateful.State[S], error)

type Collector func(ctx context.Context, persistenceId string, s stateful.State[message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error)

type OneToOneCollector[S any, OK any, OV any] func(ctx context.Context, persistenceId string, s stateful.State[S]) (*message.Message[OK, OV], error)

type OneToTwoCollector[S any, OK1 any, OV1 any, OK2 any, OV2 any] func(ctx context.Context, persistenceId string, s stateful.State[S]) (*message.Message[OK1, OV1], *message.Message[OK2, OV2], error)
