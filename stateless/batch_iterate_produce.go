package stateless

import (
	"context"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/metric"
	"github.com/hjwalt/runway/runtime"
)

// constructor
func NewProducerBatchIterateFunction(configurations ...runtime.Configuration[*ProducerBatchIterateFunction]) BatchFunction {
	batchIterateFunction := &ProducerBatchIterateFunction{}
	for _, configuration := range configurations {
		batchIterateFunction = configuration(batchIterateFunction)
	}

	batchFunction := BatchProducer{
		producer: batchIterateFunction.producer,
		next:     batchIterateFunction.Apply,
		metric:   batchIterateFunction.metric,
	}

	return batchFunction.Apply
}

// configuration

func WithBatchIterateFunctionProducer(producer message.Producer) runtime.Configuration[*ProducerBatchIterateFunction] {
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

func WithBatchIterateProducerPrometheus() runtime.Configuration[*ProducerBatchIterateFunction] {
	return func(psf *ProducerBatchIterateFunction) *ProducerBatchIterateFunction {
		psf.metric = metric.PrometheusProduce()
		return psf
	}
}

// implementation
type ProducerBatchIterateFunction struct {
	producer message.Producer
	next     SingleFunction
	metric   metric.Produce
}

func (r *ProducerBatchIterateFunction) Apply(c context.Context, ms []message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {
	res := make([]message.Message[message.Bytes, message.Bytes], 0)

	for _, m := range ms {
		if mres, err := r.next(c, m); err != nil {
			return make([]message.Message[message.Bytes, message.Bytes], 0), err
		} else {
			res = append(res, mres...)
		}
	}
	return res, nil
}
