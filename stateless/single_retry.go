package stateless

import (
	"context"
	"fmt"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/runway/logger"
	"go.uber.org/zap"
)

// constructor
func NewSingleRetry(configurations ...runtime.Configuration[*SingleRetry]) StatelessBinarySingleFunction {
	singleFunction := &SingleRetry{}
	for _, configuration := range configurations {
		singleFunction = configuration(singleFunction)
	}
	return singleFunction.Apply
}

// configurations
func WithSingleRetryRuntime(retry *runtime_retry.Retry) runtime.Configuration[*SingleRetry] {
	return func(psf *SingleRetry) *SingleRetry {
		psf.retry = retry
		return psf
	}
}

func WithSingleRetryNextFunction(next StatelessBinarySingleFunction) runtime.Configuration[*SingleRetry] {
	return func(psf *SingleRetry) *SingleRetry {
		psf.next = next
		return psf
	}
}

// implementation
type SingleRetry struct {
	retry *runtime_retry.Retry
	next  StatelessBinarySingleFunction
}

func (r *SingleRetry) Apply(c context.Context, m message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {
	msgs := make([]message.Message[message.Bytes, message.Bytes], 0)
	retryErr := r.retry.Do(func(tryCount int64) error {
		// add to prometheus
		retryGauge.WithLabelValues(m.Topic, fmt.Sprintf("%d", m.Partition)).Set(float64(tryCount))
		res, err := r.next(c, m)
		if err != nil {
			logger.Warn("retrying", zap.Int64("try", tryCount), zap.Error(err))
			return err
		}
		msgs = append(msgs, res...)
		return nil
	})
	return msgs, retryErr
}
