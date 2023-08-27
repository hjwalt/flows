package stateful

import (
	"context"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/runway/structure"
)

type PersistenceIdFunction[IK any, IV any] func(context.Context, flow.Message[IK, IV]) (string, error)

type SingleFunction func(context.Context, flow.Message[structure.Bytes, structure.Bytes], State[structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], State[structure.Bytes], error)

type OneToOneFunction[S any, IK any, IV any, OK any, OV any] func(context.Context, flow.Message[IK, IV], State[S]) (*flow.Message[OK, OV], State[S], error)

type OneToTwoFunction[S any, IK any, IV any, OK1 any, OV1 any, OK2 any, OV2 any] func(context.Context, flow.Message[IK, IV], State[S]) (*flow.Message[OK1, OV1], *flow.Message[OK2, OV2], State[S], error)
