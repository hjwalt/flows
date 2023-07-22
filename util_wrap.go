package flows

import (
	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/runtime"
)

func WrapRetry(fn stateless.SingleFunction, retryConfig []runtime.Configuration[*runtime_retry.Retry]) (stateless.SingleFunction, runtime.Runtime) {
	retryRuntime := runtime_retry.NewRetry(retryConfig...)

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
