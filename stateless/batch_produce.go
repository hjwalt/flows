package stateless

import (
	"context"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/metric"
	"github.com/hjwalt/flows/runtime"
)

// constructor
func NewProducerBatchFunction(configurations ...runtime.Configuration[*BatchProducer]) BatchFunction {
	batchFunction := &BatchProducer{}
	for _, configuration := range configurations {
		batchFunction = configuration(batchFunction)
	}
	return batchFunction.Apply
}

// configuration
func WithBatchProducerRuntime(producer runtime.Producer) runtime.Configuration[*BatchProducer] {
	return func(pbf *BatchProducer) *BatchProducer {
		pbf.producer = producer
		return pbf
	}
}

func WithBatchProducerNextFunction(next BatchFunction) runtime.Configuration[*BatchProducer] {
	return func(pbf *BatchProducer) *BatchProducer {
		pbf.next = next
		return pbf
	}
}

// implementation
type BatchProducer struct {
	producer runtime.Producer
	next     BatchFunction
	metric   metric.Produce
}

func (r BatchProducer) Apply(c context.Context, m []message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {
	n, err := r.next(c, m)
	if err != nil {
		return make([]message.Message[message.Bytes, message.Bytes], 0), err
	}
	if err := r.producer.Produce(c, n); err != nil {
		return make([]message.Message[message.Bytes, message.Bytes], 0), err
	}
	if r.metric != nil {
		r.metric.MessagesProducedIncrement(m[0].Topic, m[0].Partition, int64(len(n)))
	}
	return make([]message.Message[message.Bytes, message.Bytes], 0), nil
}
