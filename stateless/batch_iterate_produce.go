package stateless

import (
	"context"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime"
)

// constructor
func NewProducerBatchIterateFunction(configurations ...runtime.Configuration[*ProducerBatchIterateFunction]) BatchFunction {
	batchIterateFunction := &ProducerBatchIterateFunction{}
	for _, configuration := range configurations {
		batchIterateFunction = configuration(batchIterateFunction)
	}

	// wrap next function with produce per iteration
	withProducerBatchIterateWrap(batchIterateFunction)

	return batchIterateFunction.Apply
}

// configuration

func WithBatchIterateFunctionProducer(producer runtime.Producer) runtime.Configuration[*ProducerBatchIterateFunction] {
	return func(pbif *ProducerBatchIterateFunction) *ProducerBatchIterateFunction {
		pbif.producer = producer
		return pbif
	}
}

func WithBatchIterateFunctionNextFunction(next SingleFunction) runtime.Configuration[*ProducerBatchIterateFunction] {
	return func(pbif *ProducerBatchIterateFunction) *ProducerBatchIterateFunction {
		pbif.next = next
		return pbif
	}
}

func withProducerBatchIterateWrap(c *ProducerBatchIterateFunction) {
	c.next = NewSingleProducer(
		WithSingleProducerRuntime(c.producer),
		WithSingleProducerNextFunction(c.next),
	)
}

// implementation
type ProducerBatchIterateFunction struct {
	producer runtime.Producer
	next     SingleFunction
}

func (r *ProducerBatchIterateFunction) Apply(c context.Context, ms []message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {
	for _, m := range ms {
		if _, err := r.next(c, m); err != nil {
			return make([]message.Message[message.Bytes, message.Bytes], 0), err
		}
	}
	return make([]message.Message[message.Bytes, message.Bytes], 0), nil
}
