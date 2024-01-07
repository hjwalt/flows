package task_executor_retry

import (
	"context"
	"errors"

	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/task"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/runtime"
	"github.com/hjwalt/runway/structure"
	"go.uber.org/zap"
)

// constructor
func New(configurations ...runtime.Configuration[*TaskRetry]) task.Executor[structure.Bytes] {
	function := &TaskRetry{}
	for _, configuration := range configurations {
		function = configuration(function)
	}
	return function.Apply
}

// configurations
func WithRetry(r *runtime_retry.Retry) runtime.Configuration[*TaskRetry] {
	return func(p *TaskRetry) *TaskRetry {
		p.retry = r
		return p
	}
}

func WithExecutor(e task.Executor[structure.Bytes]) runtime.Configuration[*TaskRetry] {
	return func(p *TaskRetry) *TaskRetry {
		p.executor = e
		return p
	}
}

// implementation
type TaskRetry struct {
	retry    *runtime_retry.Retry
	executor task.Executor[structure.Bytes]
}

func (r *TaskRetry) Apply(c context.Context, t task.Message[structure.Bytes]) error {
	retryErr := r.retry.Do(func(tryCount int64) error {
		err := r.executor(runtime_retry.SetTryCount(c, tryCount), t)
		if err != nil {
			logger.Warn("retrying", zap.Int64("try", tryCount), zap.Error(err))
			return err
		}
		return nil
	})

	if retryErr != nil {
		return errors.Join(ErrorRetryAttempt, retryErr)
	} else {
		return nil
	}
}

var ErrorRetryAttempt = errors.New("retry all attempts failed")
