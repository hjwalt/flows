package flows

import (
	"context"

	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
)

func InjectedRuntimes() []runtime.Runtime {
	allRuntimes, err := inverse.GetAll[runtime.Runtime](context.Background(), QualifierRuntime)
	if err != nil {
		panic(err)
	}
	return allRuntimes
}
