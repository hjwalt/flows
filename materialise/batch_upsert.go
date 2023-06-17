package materialise

import (
	"context"
	"errors"

	"github.com/hjwalt/flows/message"
)

type BatchUpsert[T any] struct {
	repository UpsertRepository[T]
	mapper     MapFunction[message.Bytes, message.Bytes, T]
}

func (r *BatchUpsert[T]) Apply(c context.Context, ms []message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {
	emptyResult := make([]message.Message[[]byte, []byte], 0)

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
