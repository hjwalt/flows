package flows

import (
	"context"

	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
)

const (
	QualifierRuntime = "QualifierRuntime"
)

func InjectorRuntime(qualifier string) inverse.Injector[runtime.Runtime] {
	return func(ctx context.Context) (runtime.Runtime, error) {
		return inverse.GetLast[runtime.Runtime](ctx, qualifier)
	}
}

func InjectedRuntimes() []runtime.Runtime {
	allRuntimes, err := inverse.GetAll[runtime.Runtime](context.Background(), QualifierRuntime)
	if err != nil {
		panic(err)
	}
	return allRuntimes
}
