package stateful

import (
	"context"

	"github.com/hjwalt/flows/message"
)

type PersistenceIdFunction[IK any, IV any] func(context.Context, message.Message[IK, IV]) (string, error)

type SingleFunction func(context.Context, message.Message[message.Bytes, message.Bytes], SingleState[message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], SingleState[message.Bytes], error)

type OneToOneFunction[S any, IK any, IV any, OK any, OV any] func(context.Context, message.Message[IK, IV], SingleState[S]) (*message.Message[OK, OV], SingleState[S], error)

type OneToTwoFunction[S any, IK any, IV any, OK1 any, OV1 any, OK2 any, OV2 any] func(context.Context, message.Message[IK, IV], SingleState[S]) (*message.Message[OK1, OV1], *message.Message[OK2, OV2], SingleState[S], error)
