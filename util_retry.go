package flows

import (
	"context"

	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
)

const (
	QualifierRetryConfiguration = "QualifierRetryConfiguration"
	QualifierRetry              = "QualifierRetry"
)

// Retry
func RegisterRetry(config []runtime.Configuration[*runtime_retry.Retry]) {
	inverse.RegisterInstances(QualifierRetryConfiguration, config)
	inverse.RegisterWithConfigurationOptional[*runtime_retry.Retry](
		QualifierRetry,
		QualifierRetryConfiguration,
		runtime_retry.NewRetry,
	)
	inverse.Register(QualifierRuntime, InjectorRuntime(QualifierRetry))
}

func GetRetry(ctx context.Context) (*runtime_retry.Retry, error) {
	return inverse.GetLast[*runtime_retry.Retry](ctx, QualifierRetry)
}
