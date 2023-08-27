package stateless

import (
	"context"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/metric"
	"github.com/hjwalt/runway/runtime"
	"github.com/hjwalt/runway/structure"
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

func WithBatchIterateFunctionProducer(producer flow.Producer) runtime.Configuration[*ProducerBatchIterateFunction] {
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
	producer flow.Producer
	next     SingleFunction
	metric   metric.Produce
}

func (r *ProducerBatchIterateFunction) Apply(c context.Context, ms []flow.Message[structure.Bytes, structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], error) {
	res := make([]flow.Message[structure.Bytes, structure.Bytes], 0)

	for _, m := range ms {
		if mres, err := r.next(c, m); err != nil {
			return make([]flow.Message[structure.Bytes, structure.Bytes], 0), err
		} else {
			res = append(res, mres...)
		}
	}
	return res, nil
}
