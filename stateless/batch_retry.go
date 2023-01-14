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
func NewBatchRetry(configurations ...runtime.Configuration[*BatchRetry]) StatelessBinaryBatchFunction {
	batchFunction := &BatchRetry{}
	for _, configuration := range configurations {
		batchFunction = configuration(batchFunction)
	}
	return batchFunction.Apply
}

// configurations
func WithBatchRetryRuntime(retry *runtime_retry.Retry) runtime.Configuration[*BatchRetry] {
	return func(psf *BatchRetry) *BatchRetry {
		psf.retry = retry
		return psf
	}
}

func WithBatchRetryNextFunction(next StatelessBinaryBatchFunction) runtime.Configuration[*BatchRetry] {
	return func(psf *BatchRetry) *BatchRetry {
		psf.next = next
		return psf
	}
}

// implementation
type BatchRetry struct {
	retry *runtime_retry.Retry
	next  StatelessBinaryBatchFunction
}

func (r *BatchRetry) Apply(c context.Context, m []message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {
	msgs := make([]message.Message[message.Bytes, message.Bytes], 0)
	retryErr := r.retry.Do(func(tryCount int64) error {
		// add to prometheus
		retryGauge.WithLabelValues(m[0].Topic, fmt.Sprintf("%d", m[0].Partition)).Set(float64(tryCount))
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
