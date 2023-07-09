package runtime_retry

import (
	"sync/atomic"

	"github.com/avast/retry-go"
	"github.com/hjwalt/runway/runtime"
)

// configurations
func WithRetryOption(options ...retry.Option) runtime.Configuration[*Retry] {
	return func(c *Retry) *Retry {
		options = append(options, retry.RetryIf(AlwaysTry))
		c.options = options
		return c
	}
}

// constructor
func NewRetry(configurations ...runtime.Configuration[*Retry]) *Retry {
	consumer := &Retry{}
	for _, configuration := range configurations {
		consumer = configuration(consumer)
	}
	return consumer
}

// implementation
type Retry struct {
	options []retry.Option
	stopped atomic.Bool
}

func (c *Retry) Start() error {
	c.stopped.Store(false)
	return nil
}

func (c *Retry) Stop() {
	c.stopped.Store(true)
}

func (c *Retry) Do(fnToDo func(int64) error) error {
	tryCount := int64(0)
	return retry.Do(func() error {
		if c.stopped.Load() {
			return Stopped()
		}
		tryCount += 1
		return fnToDo(tryCount)
	}, c.options...)
}
