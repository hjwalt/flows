package stateful

import (
	"context"

	"github.com/hjwalt/flows/message"
)

type StatefulBinaryPersistenceIdFunction func(context.Context, message.Message[message.Bytes, message.Bytes]) (string, error)

type StatefulBinarySingleFunction func(context.Context, message.Message[message.Bytes, message.Bytes], SingleState[message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], SingleState[message.Bytes], error)
