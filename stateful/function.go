package stateful

import (
	"context"

	"github.com/hjwalt/flows/message"
)

type PersistenceIdFunction func(context.Context, message.Message[message.Bytes, message.Bytes]) (string, error)

type SingleFunction func(context.Context, message.Message[message.Bytes, message.Bytes], SingleState[message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], SingleState[message.Bytes], error)
