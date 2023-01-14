package stateless

import (
	"context"

	"github.com/hjwalt/flows/message"
)

type StatelessBinarySingleFunction func(context.Context, message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error)

type StatelessBinaryBatchFunction func(context.Context, []message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error)
