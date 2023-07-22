package runtime_retry

import (
	"sync/atomic"
	"time"

	"github.com/avast/retry-go"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/runtime"
)

// constructor
func NewRetry(configurations ...runtime.Configuration[*Retry]) *Retry {
	consumer := &Retry{
		options: []retry.Option{
			retry.RetryIf(AlwaysTry),
			retry.Attempts(1000000), // more than a week as a default
			retry.Delay(10 * time.Millisecond),
			retry.MaxDelay(time.Second),
			retry.MaxJitter(time.Second),
			retry.DelayType(retry.BackOffDelay),
		},
		absorb: false,
	}
	for _, configuration := range configurations {
		consumer = configuration(consumer)
	}
	return consumer
}

// configurations
func WithRetryOption(options ...retry.Option) runtime.Configuration[*Retry] {
	return func(c *Retry) *Retry {
		c.options = append(c.options, options...)
		return c
	}
}

func WithAttempts(attempts uint) runtime.Configuration[*Retry] {
	return func(c *Retry) *Retry {
		c.options = append(c.options, retry.Attempts(attempts))
		return c
	}
}

func WithDelay(delay time.Duration) runtime.Configuration[*Retry] {
	return func(c *Retry) *Retry {
		c.options = append(c.options, retry.Delay(delay))
		return c
	}
}

func WithMaxDelay(delay time.Duration) runtime.Configuration[*Retry] {
	return func(c *Retry) *Retry {
		c.options = append(c.options, retry.MaxDelay(delay))
		return c
	}
}

func WithAbsorbError(absorb bool) runtime.Configuration[*Retry] {
	return func(c *Retry) *Retry {
		c.absorb = absorb
		return c
	}
}

// implementation
type Retry struct {
	options []retry.Option
	absorb  bool
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

	err := retry.Do(func() error {
		if c.stopped.Load() {
			return Stopped()
		}
		tryCount += 1
		return fnToDo(tryCount)
	}, c.options...)

	if c.absorb && err != nil {
		logger.ErrorErr("absorbing retry error", err)
		return nil
	}

	return err
}
