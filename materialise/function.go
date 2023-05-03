package materialise

import (
	"context"

	"github.com/hjwalt/flows/message"
)

type UpsertRepository[T any] interface {
	Upsert(context.Context, []T) error
}

type MapFunction[T any] func(context.Context, message.Message[message.Bytes, message.Bytes]) ([]T, error)
