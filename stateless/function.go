package stateless

import (
	"context"

	"github.com/hjwalt/flows/message"
)

type SingleFunction func(context.Context, message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error)

type BatchFunction func(context.Context, []message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error)

type OneToOneFunction[IK any, IV any, OK any, OV any] func(context.Context, message.Message[IK, IV]) (*message.Message[OK, OV], error)

type OneToOneExplodeFunction[IK any, IV any, OK any, OV any] func(context.Context, message.Message[IK, IV]) ([]message.Message[OK, OV], error)

type OneToTwoFunction[IK any, IV any, OK1 any, OV1 any, OK2 any, OV2 any] func(context.Context, message.Message[IK, IV]) (*message.Message[OK1, OV1], *message.Message[OK2, OV2], error)
