package flows

import (
	"context"

	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
)

const (
	QualifierRetry = "QualifierRetry"
)

// Retry
func RegisterRetry(
	container inverse.Container,
	configs []runtime.Configuration[*runtime_retry.Retry],
) {

	resolver := runtime.NewResolver[*runtime_retry.Retry, *runtime_retry.Retry](
		QualifierRetry,
		container,
		false,
		runtime_retry.NewRetry,
	)

	for _, config := range configs {
		resolver.AddConfigVal(config)
	}

	resolver.Register()

	RegisterRuntime(QualifierRetry, container)
}

func GetRetry(ctx context.Context, ci inverse.Container) (*runtime_retry.Retry, error) {
	return inverse.GenericGetLast[*runtime_retry.Retry](ci, ctx, QualifierRetry)
}
