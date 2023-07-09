package flows

import (
	"time"

	"github.com/avast/retry-go"
	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/runtime"
)

func WrapRetry(fn stateless.SingleFunction) (stateless.SingleFunction, runtime.Runtime) {
	retryRuntime := runtime_retry.NewRetry(
		runtime_retry.WithRetryOption(
			retry.Attempts(1000000),
			retry.Delay(10*time.Millisecond),
			retry.MaxDelay(time.Second),
			retry.MaxJitter(time.Second),
			retry.DelayType(retry.BackOffDelay),
		),
	)
	wrappedFunction := stateless.NewSingleRetry(
		stateless.WithSingleRetryRuntime(retryRuntime),
		stateless.WithSingleRetryNextFunction(fn),
		stateless.WithSingleRetryPrometheus(),
	)
	return wrappedFunction, retryRuntime
}

func WrapSingleProduce(fn stateless.SingleFunction, producer message.Producer) stateless.SingleFunction {
	wrappedFunction := stateless.NewSingleProducer(
		stateless.WithSingleProducerNextFunction(fn),
		stateless.WithSingleProducerRuntime(producer),
		stateless.WithSingleProducerPrometheus(),
	)
	return wrappedFunction
}
