package flows

import (
	"context"

	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
)

const (
	QualifierRuntime = "QualifierRuntime"
)

func RegisterRuntime(qualifier string) {
	inverse.Register(QualifierRuntime, InjectorRuntime(qualifier))
}

func InjectorRuntime(qualifier string) inverse.Injector[runtime.Runtime] {
	return func(ctx context.Context) (runtime.Runtime, error) {
		return inverse.GetLast[runtime.Runtime](ctx, qualifier)
	}
}
