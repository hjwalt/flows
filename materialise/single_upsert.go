package materialise

import (
	"context"
	"errors"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/flows/stateless"
)

// constructor
func NewSingleUpsert[T any](configurations ...runtime.Configuration[*SingleUpsert[T]]) stateless.SingleFunction {
	r := &SingleUpsert[T]{}

	for _, configuration := range configurations {
		r = configuration(r)
	}
	return r.Apply
}

// configuration
func WithSingleUpsertRepository[T any](repository UpsertRepository[T]) runtime.Configuration[*SingleUpsert[T]] {
	return func(c *SingleUpsert[T]) *SingleUpsert[T] {
		c.repository = repository
		return c
	}
}

func WithSingleUpsertMapFunction[T any](mapper MapFunction[message.Bytes, message.Bytes, T]) runtime.Configuration[*SingleUpsert[T]] {
	return func(c *SingleUpsert[T]) *SingleUpsert[T] {
		c.mapper = mapper
		return c
	}
}

// implementation
type SingleUpsert[T any] struct {
	repository UpsertRepository[T]
	mapper     MapFunction[message.Bytes, message.Bytes, T]
}

func (r *SingleUpsert[T]) Apply(c context.Context, m message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {
	emptyResult := make([]message.Message[[]byte, []byte], 0)

	// map
	mapped, mapErr := r.mapper(c, m)
	if mapErr != nil {
		return emptyResult, errors.Join(ErrorUpsertMapper, mapErr)
	}

	// upsert
	upsertErr := r.repository.Upsert(c, mapped)
	if upsertErr != nil {
		return emptyResult, errors.Join(ErrorUpsertRepository, upsertErr)
	}

	return emptyResult, nil
}
