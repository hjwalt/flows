package stateless

import (
	"context"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/metric"
	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/runway/logger"
	"go.uber.org/zap"
)

// constructor
func NewBatchRetry(configurations ...runtime.Configuration[*BatchRetry]) BatchFunction {
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

func WithBatchRetryNextFunction(next BatchFunction) runtime.Configuration[*BatchRetry] {
	return func(psf *BatchRetry) *BatchRetry {
		psf.next = next
		return psf
	}
}

func WithBatchRetryPrometheus() runtime.Configuration[*SingleRetry] {
	return func(sr *SingleRetry) *SingleRetry {
		sr.metric = metric.PrometheusRetry()
		return sr
	}
}

// implementation
type BatchRetry struct {
	retry  *runtime_retry.Retry
	next   BatchFunction
	metric metric.Retry
}

func (r *BatchRetry) Apply(c context.Context, m []message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {
	msgs := make([]message.Message[message.Bytes, message.Bytes], 0)
	retryErr := r.retry.Do(func(tryCount int64) error {
		if r.metric != nil {
			r.metric.RetryCount(m[0].Topic, m[0].Partition, tryCount)
		}
		res, err := r.next(runtime_retry.SetTryCount(c, tryCount), m)
		if err != nil {
			logger.Warn("retrying", zap.Int64("try", tryCount), zap.Error(err))
			return err
		}
		msgs = append(msgs, res...)
		return nil
	})
	return msgs, retryErr
}
