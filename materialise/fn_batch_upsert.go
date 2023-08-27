package materialise

import (
	"context"
	"errors"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/runtime"
	"github.com/hjwalt/runway/structure"
)

// constructor
func NewBatchUpsert[T any](configurations ...runtime.Configuration[*BatchUpsert[T]]) stateless.BatchFunction {
	r := &BatchUpsert[T]{}

	for _, configuration := range configurations {
		r = configuration(r)
	}
	return r.Apply
}

// configuration
func WithBatchUpsertRepository[T any](repository UpsertRepository[T]) runtime.Configuration[*BatchUpsert[T]] {
	return func(c *BatchUpsert[T]) *BatchUpsert[T] {
		c.repository = repository
		return c
	}
}

func WithBatchUpsertMapFunction[T any](mapper MapFunction[structure.Bytes, structure.Bytes, T]) runtime.Configuration[*BatchUpsert[T]] {
	return func(c *BatchUpsert[T]) *BatchUpsert[T] {
		c.mapper = mapper
		return c
	}
}

type BatchUpsert[T any] struct {
	repository UpsertRepository[T]
	mapper     MapFunction[structure.Bytes, structure.Bytes, T]
}

func (r *BatchUpsert[T]) Apply(c context.Context, ms []flow.Message[structure.Bytes, structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], error) {
	emptyResult := make([]flow.Message[[]byte, []byte], 0)

	entities := make([]T, 0)

	// map
	for _, m := range ms {
		mapped, mapErr := r.mapper(c, m)
		if mapErr != nil {
			return emptyResult, errors.Join(ErrorUpsertMapper, mapErr)
		}

		entities = append(entities, mapped...)
	}

	// upsert
	upsertErr := r.repository.Upsert(c, entities)
	if upsertErr != nil {
		return emptyResult, errors.Join(ErrorUpsertRepository, upsertErr)
	}

	return emptyResult, nil
}
