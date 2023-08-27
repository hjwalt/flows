package materialise

import (
	"context"

	"github.com/hjwalt/flows/flow"
)

type UpsertRepository[T any] interface {
	Upsert(context.Context, []T) error
}

type MapFunction[K any, V any, T any] func(context.Context, flow.Message[K, V]) ([]T, error)
