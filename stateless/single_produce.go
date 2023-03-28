package stateless

import (
	"context"
	"fmt"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/runway/logger"
)

// constructor
func NewSingleProducer(configurations ...runtime.Configuration[*SingleProducer]) SingleFunction {
	singleFunction := &SingleProducer{}
	for _, configuration := range configurations {
		singleFunction = configuration(singleFunction)
	}
	return singleFunction.Apply
}

// configurations
func WithSingleProducerRuntime(producer runtime.Producer) runtime.Configuration[*SingleProducer] {
	return func(psf *SingleProducer) *SingleProducer {
		psf.producer = producer
		return psf
	}
}

func WithSingleProducerNextFunction(next SingleFunction) runtime.Configuration[*SingleProducer] {
	return func(psf *SingleProducer) *SingleProducer {
		psf.next = next
		return psf
	}
}

// implementation
type SingleProducer struct {
	producer runtime.Producer
	next     SingleFunction
}

func (r SingleProducer) Apply(c context.Context, m message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {
	if r.next == nil {
		logger.Warn("missing next function for producer single function")
		return make([]message.Message[[]byte, []byte], 0), nil
	}

	n, err := r.next(c, m)
	if err != nil {
		return make([]message.Message[message.Bytes, message.Bytes], 0), err
	}
	if err := r.producer.Produce(c, n); err != nil {
		return make([]message.Message[message.Bytes, message.Bytes], 0), err
	}
	// add to prometheus
	messageProducedCounter.WithLabelValues(m.Topic, fmt.Sprintf("%d", m.Partition)).Add(float64(len(n)))
	return make([]message.Message[message.Bytes, message.Bytes], 0), nil
}
