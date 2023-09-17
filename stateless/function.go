package stateless

import (
	"context"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/runway/structure"
)

type SingleFunction func(context.Context, flow.Message[structure.Bytes, structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], error)

type BatchFunction func(context.Context, []flow.Message[structure.Bytes, structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], error)

type OneToOneFunction[IK any, IV any, OK any, OV any] func(context.Context, flow.Message[IK, IV]) (*flow.Message[OK, OV], error)

type OneToOneExplodeFunction[IK any, IV any, OK any, OV any] func(context.Context, flow.Message[IK, IV]) ([]flow.Message[OK, OV], error)

type OneToTwoFunction[IK any, IV any, OK1 any, OV1 any, OK2 any, OV2 any] func(context.Context, flow.Message[IK, IV]) (*flow.Message[OK1, OV1], *flow.Message[OK2, OV2], error)

type ErrorHandlerFunction func(context.Context, flow.Message[structure.Bytes, structure.Bytes], error) ([]flow.Message[structure.Bytes, structure.Bytes], error)
