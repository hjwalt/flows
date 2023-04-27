package stateless

import (
	"context"
	"errors"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/metric"
	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/runway/logger"
	"go.uber.org/zap"
)

// constructor
func NewSingleRetry(configurations ...runtime.Configuration[*SingleRetry]) SingleFunction {
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

func WithSingleRetryNextFunction(next SingleFunction) runtime.Configuration[*SingleRetry] {
	return func(psf *SingleRetry) *SingleRetry {
		psf.next = next
		return psf
	}
}

func WithSingleRetryPrometheus() runtime.Configuration[*SingleRetry] {
	return func(sr *SingleRetry) *SingleRetry {
		sr.metric = metric.PrometheusRetry()
		return sr
	}
}

// implementation
type SingleRetry struct {
	retry  *runtime_retry.Retry
	next   SingleFunction
	metric metric.Retry
}

func (r *SingleRetry) Apply(c context.Context, m message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {
	msgs := make([]message.Message[message.Bytes, message.Bytes], 0)
	retryErr := r.retry.Do(func(tryCount int64) error {
		if r.metric != nil {
			r.metric.RetryCount(m.Topic, m.Partition, tryCount)
		}
		res, err := r.next(runtime_retry.SetTryCount(c, tryCount), m)
		if err != nil {
			logger.Warn("retrying", zap.Int64("try", tryCount), zap.Error(err))
			return err
		}
		msgs = append(msgs, res...)
		return nil
	})

	if retryErr != nil {
		return msgs, errors.Join(ErrorRetryAttempt, retryErr)
	} else {
		return msgs, nil
	}
}
