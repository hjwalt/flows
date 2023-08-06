package flows

import (
	"context"
	"errors"

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
	inverse.Register(QualifierRetry, InjectorRetry)
	inverse.Register(QualifierRuntime, InjectorRuntime(QualifierRetry))
}

func InjectorRetry(ctx context.Context) (*runtime_retry.Retry, error) {
	configurations, getConfigurationError := inverse.GetAll[runtime.Configuration[*runtime_retry.Retry]](ctx, QualifierRetryConfiguration)
	if getConfigurationError != nil && !errors.Is(getConfigurationError, inverse.ErrNotInjected) {
		return nil, getConfigurationError
	}
	return runtime_retry.NewRetry(configurations...), nil
}

func GetRetry(ctx context.Context) (*runtime_retry.Retry, error) {
	return inverse.GetLast[*runtime_retry.Retry](ctx, QualifierRetry)
}
